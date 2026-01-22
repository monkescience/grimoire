package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

// guidanceInput is the input for the guidance tool.
type guidanceInput struct {
	Name  string   `json:"name,omitempty"  jsonschema:"Name of the guidance to load"`
	Names []string `json:"names,omitempty" jsonschema:"Multiple names to load in batch"`
}

func (s *Server) registerGuidance() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "guidance",
		Description: grimoire.BuildGuidanceDescription(s.store),
	}, s.handleGuidance)
}

func (s *Server) handleGuidance(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input guidanceInput,
) (*mcp.CallToolResult, any, error) {
	// Handle single name lookup
	if input.Name != "" {
		return s.handleGuidanceByName([]string{input.Name})
	}

	// Handle batch name lookup
	if len(input.Names) > 0 {
		return s.handleGuidanceByName(input.Names)
	}

	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: "provide name or names parameter"},
		},
	}, nil, nil
}

func (s *Server) handleGuidanceByName(names []string) (*mcp.CallToolResult, any, error) {
	slog.Debug("loading guidance by name", "names", names)

	var (
		entries  []*grimoire.Entry
		notFound []string
	)

	for _, name := range names {
		found := false

		for _, typ := range []grimoire.Type{grimoire.TypeSkill, grimoire.TypeRule} {
			entry, err := s.store.Get(typ, name)
			if err == nil {
				entries = append(entries, entry)
				found = true

				break
			}
		}

		if !found {
			notFound = append(notFound, name)
		}
	}

	if len(entries) == 0 {
		slog.Warn("guidance not found", "names", notFound)

		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("guidance not found: %v", notFound)},
			},
		}, nil, nil
	}

	slog.Debug("guidance loaded", "count", len(entries), "not_found", notFound)

	result := formatEntries(entries)
	if len(notFound) > 0 {
		result += fmt.Sprintf("\n\n---\nNot found: %v", notFound)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

// formatEntries formats multiple entries for display.
func formatEntries(entries []*grimoire.Entry) string {
	var b strings.Builder

	for i, entry := range entries {
		if i > 0 {
			b.WriteString("\n---\n\n")
		}

		fmt.Fprintf(&b, "# %s\n\n", entry.Name)
		b.WriteString(entry.Body)
	}

	return b.String()
}
