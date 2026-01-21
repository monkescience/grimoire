package grimoire

import (
	"fmt"
	"strings"
)

// BuildGuidanceDescription generates the guidance tool description from store content.
// This allows the AI to see all available guidance upfront.
func BuildGuidanceDescription(s *Store) string {
	var b strings.Builder

	b.WriteString("Load guidance by name for detailed instructions.\n")

	// Skills section
	skills := s.List(TypeSkill)
	if len(skills) > 0 {
		b.WriteString("\nSKILLS (task instructions):\n")

		for _, e := range skills {
			fmt.Fprintf(&b, "- %s%s: %s\n", e.Name, e.FormatTags(), e.Description)
		}
	}

	// Rules section
	rules := s.List(TypeRule)
	if len(rules) > 0 {
		b.WriteString("\nRULES (conventions to follow):\n")

		for _, e := range rules {
			fmt.Fprintf(&b, "- %s%s: %s\n", e.Name, e.FormatTags(), e.Description)
		}
	}

	return b.String()
}

// BuildServerInstructions generates the server instructions from store content.
func BuildServerInstructions(s *Store) string {
	var b strings.Builder

	b.WriteString("Grimoire provides project-specific guidance.\n\n")
	b.WriteString("IMPORTANT: Check relevant RULES before answering questions about their topics.\n")
	b.WriteString("Rules define conventions this project follows.\n\n")

	// Rules section with trigger hints from tags
	rules := s.List(TypeRule)
	if len(rules) > 0 {
		b.WriteString("RULES (check BEFORE answering):\n")

		for _, e := range rules {
			fmt.Fprintf(&b, "- %s: %s\n", e.Name, e.Description)

			if len(e.Tags) > 0 {
				fmt.Fprintf(&b, "  Topics: %s\n", strings.Join(e.Tags, ", "))
			}
		}

		b.WriteString("\n")
	}

	// Skills section
	skills := s.List(TypeSkill)
	if len(skills) > 0 {
		b.WriteString("SKILLS (load for tasks):\n")

		for _, e := range skills {
			fmt.Fprintf(&b, "- %s: %s\n", e.Name, e.Description)
		}

		b.WriteString("\n")
	}

	b.WriteString("Usage: guidance(name: \"<name>\")\n")

	return b.String()
}
