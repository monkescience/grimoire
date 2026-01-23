package grimoire

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Type string

const (
	TypeRule        Type = "rule"
	TypeSkill       Type = "skill"
	TypeInstruction Type = "instruction"
	TypeAgent       Type = "agent"
)

func (t Type) Valid() bool {
	switch t {
	case TypeRule, TypeSkill, TypeInstruction, TypeAgent:
		return true
	default:
		return false
	}
}

type Argument struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}

type Entry struct {
	Name string `yaml:"-"`
	Type Type   `yaml:"type"`

	// Description explains what this entry does and when to use it.
	// For skills, this should be detailed (up to 1024 chars) to help agents
	// understand when to activate the skill. Follows Agent Skills spec.
	Description string `yaml:"description"`

	// Globs are file patterns that trigger this entry (e.g., "*.go").
	// Used primarily by rules.
	Globs []string `yaml:"globs"`

	// Order controls the injection order for instructions (lower = earlier).
	Order int `yaml:"order"`

	// Arguments defines parameters that skills can accept for templating.
	Arguments []Argument `yaml:"arguments"`

	// Agents references agent names that this skill can delegate to.
	Agents []string `yaml:"agents"`

	Body string `yaml:"-"`
}

func (e *Entry) FormatGlobs() string {
	if len(e.Globs) == 0 {
		return ""
	}

	return " (" + strings.Join(e.Globs, ", ") + ")"
}

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
