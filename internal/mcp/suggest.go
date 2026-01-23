package mcp

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type suggestInput struct {
	Task   string   `json:"task,omitempty"   jsonschema:"Task description to find matching skills"`
	Files  []string `json:"files,omitempty"  jsonschema:"File paths to match against rule globs"`
	Topics []string `json:"topics,omitempty" jsonschema:"Keywords to match against rule descriptions"`
}

func (s *Server) registerSuggest() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name: "suggest",
		Description: `Suggest relevant guidance based on context.

- task: Find skills by task description (e.g., "commit", "review code")
- files: Find rules by file patterns (e.g., ["main.go"])
- topics: Find rules by keywords in description (e.g., ["error-handling"])

Returns matching entries. Use the guidance tool to load full content.`,
	}, s.handleSuggest)
}

func (s *Server) handleSuggest(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input suggestInput,
) (*mcp.CallToolResult, any, error) {
	slog.DebugContext(ctx, "suggesting guidance",
		slog.String("task", input.Task),
		slog.Any("files", input.Files),
		slog.Any("topics", input.Topics))

	if input.Task != "" {
		entries := s.store.FindByTask(input.Task)

		slog.DebugContext(ctx, "task-based suggestion completed",
			slog.Int("results", len(entries)))

		return s.entrySummaryResult(ctx, entries), nil, nil
	}

	if len(input.Files) > 0 {
		entries := s.store.FindByGlobs(input.Files)

		slog.DebugContext(ctx, "file-based suggestion completed",
			slog.Int("results", len(entries)))

		return s.entrySummaryResult(ctx, entries), nil, nil
	}

	if len(input.Topics) > 0 {
		entries := s.store.FindByTopics(input.Topics)

		slog.DebugContext(ctx, "topic-based suggestion completed",
			slog.Int("results", len(entries)))

		return s.entrySummaryResult(ctx, entries), nil, nil
	}

	return errorResultMsg("provide task, files, or topics parameter"), nil, nil
}
