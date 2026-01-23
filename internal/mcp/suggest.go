package mcp

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type suggestInput struct {
	Files  []string `json:"files,omitempty"  jsonschema:"File paths to match against rule globs"`
	Topics []string `json:"topics,omitempty" jsonschema:"Topics/tags to match against rules"`
}

func (s *Server) registerSuggest() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name: "suggest",
		Description: `Suggest relevant rules based on context.

Use this tool to discover rules that apply to your current work:
- files: Find rules by file patterns (e.g., ["main.go"] finds Go rules)
- topics: Find rules by topic tags (e.g., ["error-handling"])

Returns matching rules with their descriptions. Use the guidance tool to load full content.`,
	}, s.handleSuggest)
}

func (s *Server) handleSuggest(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input suggestInput,
) (*mcp.CallToolResult, any, error) {
	slog.DebugContext(ctx, "suggesting rules",
		slog.Any("files", input.Files),
		slog.Any("topics", input.Topics))

	// Handle file-based suggestions
	if len(input.Files) > 0 {
		entries := s.store.FindByGlobs(input.Files)

		slog.DebugContext(ctx, "file-based suggestion completed",
			slog.Int("results", len(entries)))

		return s.entrySummaryResult(ctx, entries), nil, nil
	}

	// Handle topic-based suggestions
	if len(input.Topics) > 0 {
		entries := s.store.FindByTopics(input.Topics)

		slog.DebugContext(ctx, "topic-based suggestion completed",
			slog.Int("results", len(entries)))

		return s.entrySummaryResult(ctx, entries), nil, nil
	}

	return errorResultMsg("provide files or topics parameter"), nil, nil
}
