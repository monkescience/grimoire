package content

// Type represents the kind of content.
type Type string

const (
	TypeRule   Type = "rule"
	TypePrompt Type = "prompt"
	TypeSkill  Type = "skill"
)

// Valid returns true if the type is a known content type.
func (t Type) Valid() bool {
	switch t {
	case TypeRule, TypePrompt, TypeSkill:
		return true
	default:
		return false
	}
}

// Entry represents a piece of content with metadata and body.
type Entry struct {
	// Name is the unique identifier for this entry.
	Name string `yaml:"name"`

	// Type indicates what kind of content this is.
	Type Type `yaml:"type"`

	// Title is the human-readable title.
	Title string `yaml:"title"`

	// Description is a short summary of the content.
	Description string `yaml:"description"`

	// Tags for categorization and search.
	Tags []string `yaml:"tags"`

	// Arguments for prompt templates (only used for prompts).
	Arguments []Argument `yaml:"arguments,omitempty"`

	// Body is the main content (markdown).
	Body string `yaml:"-"`
}

// Argument defines a parameter for prompt templates.
type Argument struct {
	// Name of the argument.
	Name string `yaml:"name"`

	// Description explains what this argument is for.
	Description string `yaml:"description"`

	// Required indicates if this argument must be provided.
	Required bool `yaml:"required"`
}
