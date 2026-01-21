package server

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/content"
	"github.com/monke/grimoire/internal/store"
)

// guidanceInput is the input for the guidance tool.
type guidanceInput struct {
	Name string `json:"name" jsonschema:"Name of the guidance to load"`
}

// buildGuidanceDescription generates the guidance tool description from store content.
// This allows the AI to see all available guidance upfront.
func buildGuidanceDescription(s *store.Store) string {
	var b strings.Builder

	b.WriteString("Load guidance by name for detailed instructions.\n")

	// Skills section
	skills := s.List(content.TypeSkill)
	if len(skills) > 0 {
		b.WriteString("\nSKILLS (task instructions):\n")

		for _, e := range skills {
			fmt.Fprintf(&b, "- %s: %s\n", e.Name, e.Description)
		}
	}

	// Rules section
	rules := s.List(content.TypeRule)
	if len(rules) > 0 {
		b.WriteString("\nRULES (conventions to follow):\n")

		for _, e := range rules {
			fmt.Fprintf(&b, "- %s: %s\n", e.Name, e.Description)
		}
	}

	// Prompts section
	prompts := s.List(content.TypePrompt)
	if len(prompts) > 0 {
		b.WriteString("\nPROMPTS (templates):\n")

		for _, e := range prompts {
			fmt.Fprintf(&b, "- %s: %s\n", e.Name, e.Description)
		}
	}

	return b.String()
}

// buildServerInstructions generates the server instructions from store content.
func buildServerInstructions(s *store.Store) string {
	var b strings.Builder

	b.WriteString("Grimoire provides guidance and knowledge. ")
	b.WriteString("Use the guidance tool to load relevant content when users need help.\n\n")

	// Collect all entry names for examples
	skills := s.List(content.TypeSkill)
	rules := s.List(content.TypeRule)
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

func (s *Server) registerGuidance() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "guidance",
		Description: buildGuidanceDescription(s.store),
	}, s.handleGuidance)
}

func (s *Server) handleGuidance(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input guidanceInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("loading guidance", "name", input.Name)

	// Try each type until we find a match
	for _, typ := range []content.Type{content.TypeSkill, content.TypeRule, content.TypePrompt} {
		entry, err := s.store.Get(typ, input.Name)
		if err == nil {
			slog.Debug("guidance loaded", "name", input.Name, "type", typ)

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: entry.Body},
				},
			}, nil, nil
		}
	}

	slog.Warn("guidance not found", "name", input.Name)

	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("guidance %q not found", input.Name)},
		},
	}, nil, nil
}
