package grimoire

import (
	"cmp"
	"fmt"
	"io/fs"
	"slices"
	"strings"
)

// Store provides access to grimoire content.
type Store struct {
	entries map[Type]map[string]*Entry
}

// New creates a store by loading all markdown files from the given filesystem.
// Files are scanned recursively and the type is determined from frontmatter.
func New(fsys fs.FS) (*Store, error) {
	s := &Store{
		entries: map[Type]map[string]*Entry{
			TypeRule:  {},
			TypeSkill: {},
		},
	}

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		entry, err := parseMarkdown(data)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		// Validate entry type
		if !entry.Type.Valid() {
			return fmt.Errorf("parsing %s: %w: %q", path, ErrInvalidType, entry.Type)
		}

		// Use filename (without extension) as name if not specified
		if entry.Name == "" {
			name := strings.TrimSuffix(path, ".md")
			if idx := strings.LastIndex(name, "/"); idx != -1 {
				name = name[idx+1:]
			}

			entry.Name = name
		}

		if _, exists := s.entries[entry.Type]; !exists {
			s.entries[entry.Type] = make(map[string]*Entry)
		}

		s.entries[entry.Type][entry.Name] = entry

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking filesystem: %w", err)
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

	slices.SortFunc(result, func(a, b *Entry) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return result
}

// Search finds entries matching the query across all types, sorted by name.
func (s *Store) Search(query string) []*Entry {
	query = strings.ToLower(query)

	var results []*Entry

	for _, entries := range s.entries {
		for _, entry := range entries {
			if matches(entry, query) {
				results = append(results, entry)
			}
		}
	}

	slices.SortFunc(results, func(a, b *Entry) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return results
}

// matches checks if an entry matches the search query.
func matches(entry *Entry, query string) bool {
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
