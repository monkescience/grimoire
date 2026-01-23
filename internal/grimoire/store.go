package grimoire

import (
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const minWordLength = 3

type Store struct {
	entries map[Type]map[string]*Entry
}

// New creates a store by loading content according to the provided config.
// If builtinFS is provided and config enables builtin, embedded content is loaded.
// External paths from config are also loaded.
// Returns an error if duplicate entries are found or if a source path doesn't exist.
func New(cfg *Config, builtinFS fs.FS) (*Store, error) {
	s := &Store{
		entries: map[Type]map[string]*Entry{
			TypeRule:        {},
			TypeSkill:       {},
			TypeInstruction: {},
			TypeAgent:       {},
		},
	}

	// Load external paths first (higher priority for error messages)
	for _, path := range cfg.Sources.Paths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("source %q: %w", path, ErrSourceNotFound)
			}

			return nil, fmt.Errorf("source %q: %w", path, err)
		}

		if !info.IsDir() {
			return nil, fmt.Errorf("source %q: %w", path, ErrNotDirectory)
		}

		err = s.loadFromFS(os.DirFS(path), path, cfg)
		if err != nil {
			return nil, err
		}
	}

	if cfg.BuiltinEnabled() && builtinFS != nil {
		err := s.loadFromFS(builtinFS, "builtin", cfg)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store) Get(typ Type, name string) (*Entry, error) {
	if name == "" {
		return nil, fmt.Errorf("%s: %w", typ, ErrNameEmpty)
	}

	if entries, ok := s.entries[typ]; ok {
		if entry, ok := entries[name]; ok {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("%s %q: %w", typ, name, ErrNotFound)
}

func (s *Store) List(typ Type) []*Entry {
	entries, ok := s.entries[typ]
	if !ok {
		return nil
	}

	result := make([]*Entry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry)
	}

	sortEntriesByName(result)

	return result
}

func (s *Store) Search(query string) []*Entry {
	query = strings.ToLower(query)

	var results []*Entry

	for _, entries := range s.entries {
		for _, entry := range entries {
			if matchesQuery(entry, query) {
				results = append(results, entry)
			}
		}
	}

	sortEntriesByName(results)

	return results
}

// FindByTopics returns all rules whose description matches any of the given topics.
// Matching is case-insensitive substring match against description.
func (s *Store) FindByTopics(topics []string) []*Entry {
	if len(topics) == 0 {
		return nil
	}

	// Normalize topics for comparison
	normalizedTopics := make([]string, len(topics))
	for i, t := range topics {
		normalizedTopics[i] = strings.ToLower(t)
	}

	var results []*Entry

	// Only search rules (not skills) for topic matching
	for _, entry := range s.entries[TypeRule] {
		if matchesTopics(entry, normalizedTopics) {
			results = append(results, entry)
		}
	}

	sortEntriesByName(results)

	return results
}

func (s *Store) FindByGlobs(files []string) []*Entry {
	if len(files) == 0 {
		return nil
	}

	var results []*Entry

	for _, entry := range s.entries[TypeRule] {
		if matchesGlob(entry, files) {
			results = append(results, entry)
		}
	}

	sortEntriesByName(results)

	return results
}

// FindByTask returns all skills whose description matches the given task.
// Matching is case-insensitive and checks if task keywords appear in description.
func (s *Store) FindByTask(task string) []*Entry {
	if task == "" {
		return nil
	}

	task = strings.ToLower(task)

	var results []*Entry

	for _, entry := range s.entries[TypeSkill] {
		if matchesTask(entry, task) {
			results = append(results, entry)
		}
	}

	sortEntriesByName(results)

	return results
}

// loadFromFS loads entries from a filesystem into the store.
// sourceName is used for error messages to identify the source.
func (s *Store) loadFromFS(fsys fs.FS, sourceName string, cfg *Config) error {
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, readErr := fs.ReadFile(fsys, path)
		if readErr != nil {
			return fmt.Errorf("reading %s: %w", path, readErr)
		}

		entry, parseErr := parseMarkdown(data)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", path, parseErr)
		}

		// Validate entry type
		if !entry.Type.Valid() {
			return fmt.Errorf("parsing %s: %w: %q", path, ErrInvalidType, entry.Type)
		}

		// Validate entry (globs, etc.)
		validateErr := entry.Validate()
		if validateErr != nil {
			return fmt.Errorf("validating %s: %w", path, validateErr)
		}

		// Derive name from path, stripping type prefix if present
		entry.Name = deriveName(path, entry.Type)

		// Check if entry is allowed by filter
		filter := cfg.FilterForType(entry.Type)
		if !filter.IsAllowed(entry.Name) {
			return nil // Skip filtered entries
		}

		// Check for duplicates
		if _, exists := s.entries[entry.Type]; !exists {
			s.entries[entry.Type] = make(map[string]*Entry)
		}

		if _, exists := s.entries[entry.Type][entry.Name]; exists {
			return fmt.Errorf("%s %q from %s: %w (already loaded)", entry.Type, entry.Name, sourceName, ErrDuplicate)
		}

		s.entries[entry.Type][entry.Name] = entry

		return nil
	})
	if err != nil {
		return fmt.Errorf("loading %s: %w", sourceName, err)
	}

	return nil
}

// deriveName extracts the entry name from the file path.
// It removes the .md extension and strips type-based prefixes (e.g., "rules/", "skills/").
// If no recognized prefix is found, the full relative path is preserved (minus extension).
// Examples:
//   - "rules/go/error-assignment.md" -> "go/error-assignment"
//   - "skills/refactor.md" -> "refactor"
//   - "custom/go/my-rule.md" -> "custom/go/my-rule"
//   - "my-rule.md" -> "my-rule"
func deriveName(path string, typ Type) string {
	name := strings.TrimSuffix(path, ".md")

	prefixes := []string{
		string(typ) + "s/", // "rules/", "skills/"
		string(typ) + "/",  // "rule/", "skill/" (alternative)
	}

	for _, prefix := range prefixes {
		if after, found := strings.CutPrefix(name, prefix); found {
			return after
		}
	}

	return name
}

func sortEntriesByName(entries []*Entry) {
	slices.SortFunc(entries, func(a, b *Entry) int {
		return cmp.Compare(a.Name, b.Name)
	})
}

func matchesQuery(entry *Entry, query string) bool {
	if strings.Contains(strings.ToLower(entry.Name), query) {
		return true
	}

	if strings.Contains(strings.ToLower(entry.Description), query) {
		return true
	}

	return strings.Contains(strings.ToLower(entry.Body), query)
}

// matchesTopics checks if any of the topics appear in the entry's description.
func matchesTopics(entry *Entry, topics []string) bool {
	desc := strings.ToLower(entry.Description)

	for _, topic := range topics {
		if strings.Contains(desc, topic) {
			return true
		}
	}

	return false
}

func matchesGlob(entry *Entry, files []string) bool {
	for _, pattern := range entry.Globs {
		for _, file := range files {
			matched, err := filepath.Match(pattern, file)
			if err == nil && matched {
				return true
			}

			matched, err = filepath.Match(pattern, filepath.Base(file))
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}

// matchesTask checks if the task description matches the skill's description.
// It tokenizes the task into words and checks if any significant word
// appears in the skill's description.
func matchesTask(entry *Entry, task string) bool {
	desc := strings.ToLower(entry.Description)

	// Check if the full task appears as substring
	if strings.Contains(desc, task) {
		return true
	}

	// Tokenize task into words and check for matches
	for word := range strings.FieldsSeq(task) {
		// Skip very short words (articles, prepositions)
		if len(word) < minWordLength {
			continue
		}

		if strings.Contains(desc, word) {
			return true
		}
	}

	return false
}
