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

func BuildAgentDescription(s *Store) string {
	var b strings.Builder

	b.WriteString("Execute agent prompts via MCP sampling. Requires client sampling support.\n")

	agents := s.List(TypeAgent)
	if len(agents) > 0 {
		b.WriteString("\nAGENTS:\n")

		for _, e := range agents {
			fmt.Fprintf(&b, "- %s: %s\n", e.Name, e.Description)
		}
	}

	return b.String()
}

// summarizeDescription returns a short summary of the description.
// Uses the first sentence or line, truncated if too long.
func summarizeDescription(desc string) string {
	const maxLen = 120

	// Find first sentence or line break
	desc = strings.TrimSpace(desc)

	if idx := strings.Index(desc, "\n"); idx > 0 && idx < maxLen {
		desc = desc[:idx]
	}

	if idx := strings.Index(desc, ". "); idx > 0 && idx < maxLen {
		desc = desc[:idx+1]
	}

	// Truncate if still too long
	if len(desc) > maxLen {
		desc = desc[:maxLen-3] + "..."
	}

	return desc
}

// BuildServerInstructions generates the server instructions.
// Includes skills, rules, and instruction entries for the AI to see upfront.
func BuildServerInstructions(s *Store) string {
	var b strings.Builder

	b.WriteString("Grimoire provides project-specific coding guidance.\n\n")

	skills := s.List(TypeSkill)
	if len(skills) > 0 {
		b.WriteString("SKILLS - Load with guidance(name) BEFORE these tasks:\n")

		for _, e := range skills {
			// Use first line of description as summary, or truncate if needed
			summary := summarizeDescription(e.Description)
			fmt.Fprintf(&b, "- %s: %s\n", e.Name, summary)
		}

		b.WriteString("\n")
	}

	rules := s.List(TypeRule)
	if len(rules) > 0 {
		b.WriteString("RULES - Apply based on description, load with guidance() if you need examples:\n")

		for _, e := range rules {
			fmt.Fprintf(&b, "- %s%s: %s\n", e.Name, e.FormatGlobs(), e.Description)
		}

		b.WriteString("\n")
	}

	instructions := s.List(TypeInstruction)
	if len(instructions) > 0 {
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
