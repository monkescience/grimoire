package grimoire

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const frontmatterDelimiter = "---"

// parseMarkdown parses a markdown file with YAML frontmatter.
// The type is determined from the frontmatter "type" field.
func parseMarkdown(data []byte) (*Entry, error) {
	entry := &Entry{}

	if !bytes.HasPrefix(data, []byte(frontmatterDelimiter)) {
		entry.Body = string(data)

		return entry, nil
	}

	rest := data[len(frontmatterDelimiter):]
	frontmatter, body, found := bytes.Cut(rest, []byte("\n"+frontmatterDelimiter))

	if !found {
		return nil, fmt.Errorf("unclosed: %w", ErrInvalidFrontmatter)
	}

	err := yaml.Unmarshal(frontmatter, entry)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	entry.Body = strings.TrimPrefix(string(body), "\n")

	return entry, nil
}
