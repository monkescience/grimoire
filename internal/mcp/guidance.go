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
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input guidanceInput,
) (*mcp.CallToolResult, any, error) {
	if input.Name != "" {
		return s.handleGuidanceByName(ctx, []string{input.Name})
	}

	if len(input.Names) > 0 {
		return s.handleGuidanceByName(ctx, input.Names)
	}

	return errorResultMsg("provide name or names parameter"), nil, nil
}

func (s *Server) handleGuidanceByName(ctx context.Context, names []string) (*mcp.CallToolResult, any, error) {
	slog.DebugContext(ctx, "loading guidance by name", slog.Any("names", names))

	var (
		entries  []*grimoire.Entry
		notFound []string
	)

	for _, name := range names {
		found := false

		for _, typ := range []grimoire.Type{grimoire.TypeSkill, grimoire.TypeRule, grimoire.TypeAgent} {
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
		slog.WarnContext(ctx, "guidance not found", slog.Any("names", notFound))

		return errorResultMsg(fmt.Sprintf("guidance not found: %v", notFound)), nil, nil
	}

	slog.DebugContext(ctx, "guidance loaded", slog.Int("count", len(entries)), slog.Any("not_found", notFound))

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
