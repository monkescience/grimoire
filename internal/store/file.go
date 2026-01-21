package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/monke/grimoire/internal/content"
)

// FileStore loads content from the filesystem.
type FileStore struct {
	basePath string
}

// NewFileStore creates a store that reads from the given directory.
// The directory should contain subdirectories: rules/, prompts/, skills/.
func NewFileStore(basePath string) *FileStore {
	return &FileStore{basePath: basePath}
}

// Get retrieves a single entry by type and name.
func (s *FileStore) Get(typ content.Type, name string) (*content.Entry, error) {
	dir := s.typeDir(typ)
	path := filepath.Join(dir, name+".md")

	data, err := os.ReadFile(path) //nolint:gosec // path is constructed from validated input
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s %q: %w", typ, name, ErrNotFound)
		}

		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	entry, err := parseMarkdown(data, typ)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	// Use filename as name if not specified in frontmatter
	if entry.Name == "" {
		entry.Name = name
	}

	return entry, nil
}

// List returns all entries of a given type.
func (s *FileStore) List(typ content.Type) ([]*content.Entry, error) {
	dir := s.typeDir(typ)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No entries of this type
		}

		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	var result []*content.Entry

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(e.Name(), ".md")

		entry, err := s.Get(typ, name)
		if err != nil {
			continue // Skip entries that fail to parse
		}

		result = append(result, entry)
	}

	return result, nil
}

// Search finds entries matching the query across all types.
func (s *FileStore) Search(query string) ([]*content.Entry, error) {
	query = strings.ToLower(query)

	var results []*content.Entry

	for _, typ := range []content.Type{content.TypeRule, content.TypePrompt, content.TypeSkill} {
		entries, err := s.List(typ)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if s.matches(entry, query) {
				results = append(results, entry)
			}
		}
	}

	return results, nil
}

// matches checks if an entry matches the search query.
func (s *FileStore) matches(entry *content.Entry, query string) bool {
	// Check name
	if strings.Contains(strings.ToLower(entry.Name), query) {
		return true
	}

	// Check title
	if strings.Contains(strings.ToLower(entry.Title), query) {
		return true
	}

	// Check description
	if strings.Contains(strings.ToLower(entry.Description), query) {
		return true
	}

	// Check tags
	for _, tag := range entry.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	// Check body
	if strings.Contains(strings.ToLower(entry.Body), query) {
		return true
	}

	return false
}

// typeDir returns the directory path for a content type.
func (s *FileStore) typeDir(typ content.Type) string {
	return filepath.Join(s.basePath, string(typ)+"s")
}
