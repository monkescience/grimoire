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

// Store provides access to grimoire content.
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

	// Load builtin content if enabled
	if cfg.BuiltinEnabled() && builtinFS != nil {
		err := s.loadFromFS(builtinFS, "builtin", cfg)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Get retrieves a single entry by type and name.
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

// List returns all entries of a given type, sorted by name.
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

// Search finds entries matching the query across all types, sorted by name.
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

// FindByTopics returns all rules whose tags match any of the given topics.
// Matching is case-insensitive.
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
		if matchesTopic(entry, normalizedTopics) {
			results = append(results, entry)
		}
	}

	sortEntriesByName(results)

	return results
}

// FindByGlobs returns all rules whose glob patterns match any of the given file paths.
func (s *Store) FindByGlobs(files []string) []*Entry {
	if len(files) == 0 {
		return nil
	}

	var results []*Entry

	// Only search rules (not skills) for glob matching
	for _, entry := range s.entries[TypeRule] {
		if matchesGlob(entry, files) {
			results = append(results, entry)
		}
	}

	sortEntriesByName(results)

	return results
}

// FindByTriggers returns all skills whose triggers match the given task description.
func (s *Store) FindByTriggers(task string) []*Entry {
	if task == "" {
		return nil
	}

	task = strings.ToLower(task)

	var results []*Entry

	for _, entry := range s.entries[TypeSkill] {
		if matchesTrigger(entry, task) {
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
	// Remove .md extension
	name := strings.TrimSuffix(path, ".md")

	// Try to strip type-based prefix (e.g., "rules/", "skills/")
	prefixes := []string{
		string(typ) + "s/", // "rules/", "skills/"
		string(typ) + "/",  // "rule/", "skill/" (alternative)
	}

	for _, prefix := range prefixes {
		if after, found := strings.CutPrefix(name, prefix); found {
			return after
		}
	}

	// No prefix found - preserve full relative path
	return name
}

// sortEntriesByName sorts entries alphabetically by name.
func sortEntriesByName(entries []*Entry) {
	slices.SortFunc(entries, func(a, b *Entry) int {
		return cmp.Compare(a.Name, b.Name)
	})
}

// matchesQuery checks if an entry matches the search query.
func matchesQuery(entry *Entry, query string) bool {
	if strings.Contains(strings.ToLower(entry.Name), query) {
		return true
	}

	if strings.Contains(strings.ToLower(entry.Description), query) {
		return true
	}

	for _, tag := range entry.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	return strings.Contains(strings.ToLower(entry.Body), query)
}

// matchesTopic checks if an entry's tags match any of the normalized topics.
func matchesTopic(entry *Entry, topics []string) bool {
	for _, tag := range entry.Tags {
		normalizedTag := strings.ToLower(tag)

		if slices.Contains(topics, normalizedTag) {
			return true
		}
	}

	return false
}

// matchesGlob checks if any of the entry's globs match any of the given files.
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

// matchesTrigger checks if the task contains any of the entry's tags.
func matchesTrigger(entry *Entry, task string) bool {
	for _, tag := range entry.Tags {
		if strings.Contains(task, strings.ToLower(tag)) {
			return true
		}
	}

	return false
}
