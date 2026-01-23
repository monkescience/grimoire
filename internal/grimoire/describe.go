package grimoire

import (
	"cmp"
	"fmt"
	"slices"
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

	// Agents section
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
// It includes all instruction entries sorted by order, then by name.
func BuildServerInstructions(s *Store) string {
	var b strings.Builder

	b.WriteString("Grimoire provides project-specific coding guidance through skills, rules, and agents.\n\n")
	b.WriteString("SKILLS define HOW to perform tasks (commit, review, debug, refactor).\n")
	b.WriteString("→ Use suggest(task: \"...\") to find skills before starting tasks.\n")
	b.WriteString("→ Load with guidance(name: \"skill-name\").\n\n")
	b.WriteString("RULES define project conventions for file types (matched by tags/globs).\n")
	b.WriteString("→ Apply rules based on their description. Load only if you need examples.\n\n")
	b.WriteString("AGENTS are subagent prompts for specialized tasks (requires MCP sampling).\n")
	b.WriteString("→ Execute via agent(names: [...]) when client supports sampling.\n\n")
	b.WriteString("Use guidance(name: \"rule-name\") to load, search(query: \"...\") to find.\n")

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

		b.WriteString("\n---\n\n")

		for i, instr := range instructions {
			if i > 0 {
				b.WriteString("\n")
			}

			b.WriteString(instr.Body)
		}
	}

	return b.String()
}
