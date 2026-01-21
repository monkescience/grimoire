package store

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/monke/grimoire/internal/content"
)

// frontmatterDelimiter separates YAML frontmatter from markdown body.
const frontmatterDelimiter = "---"

// parseMarkdown parses a markdown file with YAML frontmatter.
func parseMarkdown(data []byte, typ content.Type) (*content.Entry, error) {
	entry := &content.Entry{Type: typ}

	// Check for frontmatter delimiter
	if !bytes.HasPrefix(data, []byte(frontmatterDelimiter)) {
		// No frontmatter, treat entire content as body
		entry.Body = string(data)

		return entry, nil
	}

	// Find the closing delimiter
	rest := data[len(frontmatterDelimiter):]
	frontmatter, body, found := bytes.Cut(rest, []byte("\n"+frontmatterDelimiter))

	if !found {
		return nil, fmt.Errorf("unclosed: %w", ErrInvalidFrontmatter)
	}

	// Parse YAML frontmatter
	err := yaml.Unmarshal(frontmatter, entry)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	// Trim leading newlines from body
	entry.Body = strings.TrimPrefix(string(body), "\n")

	return entry, nil
}
