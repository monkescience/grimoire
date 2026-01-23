package mcp

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

func (s *Server) registerInstructions() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "instructions",
		Description: "Load project instructions. Call this ONCE at the start of each session.",
	}, s.handleInstructions)
}

func (s *Server) handleInstructions(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, any, error) {
	slog.DebugContext(ctx, "loading instructions")

	instructions := s.store.List(grimoire.TypeInstruction)
	if len(instructions) == 0 {
		slog.DebugContext(ctx, "no instructions configured")

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No project instructions configured."},
			},
		}, nil, nil
	}

	slog.DebugContext(ctx, "instructions loaded", slog.Int("count", len(instructions)))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatEntries(instructions)},
		},
	}, nil, nil
}
