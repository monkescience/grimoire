package store

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/monke/grimoire/internal/content"
)

// EmbedStore loads content from an embedded filesystem.
type EmbedStore struct {
	fs     embed.FS
	prefix string
}

// NewEmbedStore creates a store that reads from an embedded filesystem.
// The filesystem should contain subdirectories: rules/, prompts/, skills/.
// The prefix is prepended to all paths (e.g., "content" if using //go:embed content).
func NewEmbedStore(efs embed.FS, prefix string) *EmbedStore {
	return &EmbedStore{fs: efs, prefix: prefix}
}

// Get retrieves a single entry by type and name.
func (s *EmbedStore) Get(typ content.Type, name string) (*content.Entry, error) {
	dir := s.typeDir(typ)
	path := dir + "/" + name + ".md"

	data, err := s.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s %q: %w", typ, name, ErrNotFound)
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
func (s *EmbedStore) List(typ content.Type) ([]*content.Entry, error) {
	dir := s.typeDir(typ)

	entries, err := fs.ReadDir(s.fs, dir)
	if err != nil {
		// No entries of this type
		return nil, nil //nolint:nilerr // missing dir is not an error
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
func (s *EmbedStore) Search(query string) ([]*content.Entry, error) {
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
func (s *EmbedStore) matches(entry *content.Entry, query string) bool {
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
func (s *EmbedStore) typeDir(typ content.Type) string {
	if s.prefix != "" {
		return s.prefix + "/" + string(typ) + "s"
	}

	return string(typ) + "s"
}
