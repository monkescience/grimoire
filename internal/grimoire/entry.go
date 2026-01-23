package grimoire

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Type represents the kind of content.
type Type string

const (
	TypeRule        Type = "rule"
	TypeSkill       Type = "skill"
	TypeInstruction Type = "instruction"
	TypeAgent       Type = "agent"
)

// Valid returns true if the type is a known content type.
func (t Type) Valid() bool {
	switch t {
	case TypeRule, TypeSkill, TypeInstruction, TypeAgent:
		return true
	default:
		return false
	}
}

// Argument describes a parameter that a skill can accept.
type Argument struct {
	// Name is the argument identifier used in templates.
	Name string `yaml:"name"`

	// Description explains what the argument is for.
	Description string `yaml:"description"`

	// Required indicates if the argument must be provided.
	Required bool `yaml:"required"`
}

// Entry represents a piece of content with metadata and body.
type Entry struct {
	// Name is the unique identifier, derived from filename.
	Name string `yaml:"-"`

	// Type indicates what kind of content this is.
	Type Type `yaml:"type"`

	// Description is a short summary of the content.
	Description string `yaml:"description"`

	// Tags for categorization and search.
	Tags []string `yaml:"tags"`

	// Globs are file patterns that trigger this entry (e.g., "*.go").
	Globs []string `yaml:"globs"`

	// Order controls the injection order for instructions (lower = earlier).
	Order int `yaml:"order"`

	// Arguments defines parameters that skills can accept for templating.
	Arguments []Argument `yaml:"arguments"`

	// MaxTokens limits response length for agents (default: 4096).
	MaxTokens int64 `yaml:"max_tokens"`

	// Agents references agent names that this skill can delegate to.
	Agents []string `yaml:"agents"`

	// Body is the main content (markdown).
	Body string `yaml:"-"`
}

// FormatTags formats tags for display in descriptions.
func (e *Entry) FormatTags() string {
	if len(e.Tags) == 0 {
		return ""
	}

	return " [" + strings.Join(e.Tags, ", ") + "]"
}

// FormatGlobs formats globs for display in descriptions.
func (e *Entry) FormatGlobs() string {
	if len(e.Globs) == 0 {
		return ""
	}

	return " (" + strings.Join(e.Globs, ", ") + ")"
}

// Validate checks the entry for errors.
func (e *Entry) Validate() error {
	for _, pattern := range e.Globs {
		_, err := filepath.Match(pattern, "")
		if err != nil {
			return fmt.Errorf("%w: %q: %w", ErrInvalidGlob, pattern, err)
		}
	}

	return nil
}

// RenderBody substitutes argument values into the body using {{argName}} syntax.
// Arguments not provided in values are replaced with empty strings.
func (e *Entry) RenderBody(values map[string]string) string {
	result := e.Body

	for _, arg := range e.Arguments {
		placeholder := "{{" + arg.Name + "}}"
		value := values[arg.Name]
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}
