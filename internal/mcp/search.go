package mcp

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// searchInput is the input for the search tool.
type searchInput struct {
	Query string `json:"query" jsonschema:"Search query"`
}

func (s *Server) registerSearch() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "search",
		Description: "Search for guidance by keyword. Returns matching skills, rules, and prompts.",
	}, s.handleSearch)
}

func (s *Server) handleSearch(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input searchInput,
) (*mcp.CallToolResult, any, error) {
	slog.DebugContext(ctx, "searching", "query", input.Query)

	entries := s.store.Search(input.Query)

	slog.DebugContext(ctx, "search completed", "query", input.Query, "results", len(entries))

	return s.entrySummaryResult(ctx, entries), nil, nil
}
