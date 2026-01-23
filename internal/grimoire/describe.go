package grimoire

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// BuildGuidanceDescription generates the guidance tool description.
// Skills and rules are listed in server instructions; this just provides usage info.
func BuildGuidanceDescription() string {
	return `Load guidance by name.

USAGE:
- guidance(name: "rule-name") - Load one
- guidance(names: ["a", "b"]) - Load multiple`
}

// BuildAgentDescription generates the agent tool description from store content.
func BuildAgentDescription(s *Store) string {
	var b strings.Builder

	b.WriteString("Execute agent prompts via MCP sampling. Requires client sampling support.\n")

	agents := s.List(TypeAgent)
	if len(agents) > 0 {
		b.WriteString("\nAGENTS:\n")

		for _, e := range agents {
			fmt.Fprintf(&b, "- %s%s: %s\n", e.Name, e.FormatTags(), e.Description)
		}
	}

	return b.String()
}

// BuildServerInstructions generates the server instructions.
// Includes skills, rules, and instruction entries for the AI to see upfront.
func BuildServerInstructions(s *Store) string {
	var b strings.Builder

	b.WriteString("Grimoire provides project-specific coding guidance.\n\n")

	// Skills section
	skills := s.List(TypeSkill)
	if len(skills) > 0 {
		b.WriteString("SKILLS - Load and follow BEFORE starting tasks:\n")

		for _, e := range skills {
			fmt.Fprintf(&b, "- %s%s: %s\n", e.Name, e.FormatTags(), e.Description)
		}

		b.WriteString("\n")
	}

	// Rules section
	rules := s.List(TypeRule)
	if len(rules) > 0 {
		b.WriteString("RULES - Apply based on description, load with guidance() if you need examples:\n")

		for _, e := range rules {
			fmt.Fprintf(&b, "- %s%s%s: %s\n", e.Name, e.FormatTags(), e.FormatGlobs(), e.Description)
		}

		b.WriteString("\n")
	}

	// Append instruction entries
	instructions := s.List(TypeInstruction)
	if len(instructions) > 0 {
		// Sort by order first, then by name
		slices.SortFunc(instructions, func(a, b *Entry) int {
			if a.Order != b.Order {
				return cmp.Compare(a.Order, b.Order)
			}

			return cmp.Compare(a.Name, b.Name)
		})

		b.WriteString("---\n\n")

		for i, instr := range instructions {
			if i > 0 {
				b.WriteString("\n")
			}

			b.WriteString(instr.Body)
		}
	}

	return b.String()
}
