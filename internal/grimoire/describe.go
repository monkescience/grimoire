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

	b.WriteString("Grimoire provides guidance and knowledge. ")
	b.WriteString("Use the guidance tool to load relevant content when users need help.\n\n")

	// Collect all entry names for examples
	skills := s.List(TypeSkill)
	rules := s.List(TypeRule)
	examples := make([]string, 0, len(skills)+len(rules))

	for _, e := range skills {
		examples = append(examples, fmt.Sprintf("For %s tasks, call guidance(name: %q)", e.Name, e.Name))
	}

	for _, e := range rules {
		examples = append(examples, fmt.Sprintf("For %s guidance, call guidance(name: %q)", e.Name, e.Name))
	}

	if len(examples) > 0 {
		b.WriteString("Examples:\n")

		for _, ex := range examples {
			fmt.Fprintf(&b, "- %s\n", ex)
		}
	}

	return b.String()
}
