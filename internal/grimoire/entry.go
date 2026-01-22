package grimoire

import "strings"

// Type represents the kind of content.
type Type string

const (
	TypeRule  Type = "rule"
	TypeSkill Type = "skill"
)

// Valid returns true if the type is a known content type.
func (t Type) Valid() bool {
	switch t {
	case TypeRule, TypeSkill:
		return true
	default:
		return false
	}
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
