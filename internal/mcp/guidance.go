package mcp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

// guidanceInput is the input for the guidance tool.
type guidanceInput struct {
	Name string `json:"name" jsonschema:"Name of the guidance to load"`
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
	slog.Debug("loading guidance", "name", input.Name)

	// Try each type until we find a match
	for _, typ := range []grimoire.Type{grimoire.TypeSkill, grimoire.TypeRule} {
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
