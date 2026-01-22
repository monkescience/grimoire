package grimoire

import (
	"fmt"
	"strings"
)

// BuildGuidanceDescription generates the guidance tool description from store content.
// This allows the AI to see all available guidance upfront.
func BuildGuidanceDescription(s *Store) string {
	var b strings.Builder

	b.WriteString("Load guidance by name.\n\n")
	b.WriteString("USAGE:\n")
	b.WriteString("- guidance(name: \"rule-name\") - Load one\n")
	b.WriteString("- guidance(names: [\"a\", \"b\"]) - Load multiple\n")

	// Skills section
	skills := s.List(TypeSkill)
	if len(skills) > 0 {
		b.WriteString("\nSKILLS:\n")

		for _, e := range skills {
			fmt.Fprintf(&b, "- %s%s: %s\n", e.Name, e.FormatTags(), e.Description)
		}
	}

	// Rules section
	rules := s.List(TypeRule)
	if len(rules) > 0 {
		b.WriteString("\nRULES:\n")

		for _, e := range rules {
			fmt.Fprintf(&b, "- %s%s%s: %s\n", e.Name, e.FormatTags(), e.FormatGlobs(), e.Description)
		}
	}

	return b.String()
}

// BuildServerInstructions generates the server instructions from store content.
func BuildServerInstructions(s *Store) string {
	var b strings.Builder

	b.WriteString("Grimoire provides project-specific guidance.\n\n")
	b.WriteString("Load relevant rules with guidance(names: [...]) before writing or reviewing code.\n")
	b.WriteString("Match rules by tags (e.g., [go]) or file patterns (e.g., (*.go)).\n")

	return b.String()
}
