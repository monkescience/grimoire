package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

// entrySummary is a lightweight representation of an entry for tool result output.
// Used by search and suggest tools to return concise entry information.
type entrySummary struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (s *Server) entrySummaryResult(ctx context.Context, entries []*grimoire.Entry) *mcp.CallToolResult {
	summaries := make([]entrySummary, len(entries))
	for i, e := range entries {
		summaries[i] = entrySummary{
			Name:        e.Name,
			Type:        string(e.Type),
			Description: e.Description,
			Tags:        e.Tags,
		}
	}

	data, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal entry summaries", slog.Any("error", err))

		return errorResult(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}
}

func (s *Server) getResourceContents(
	ctx context.Context,
	typ grimoire.Type,
	name, uri string,
) (*mcp.ReadResourceResult, error) {
	entry, err := s.store.Get(typ, name)
	if err != nil {
		slog.WarnContext(ctx, "failed to get resource contents",
			slog.String("type", string(typ)), slog.String("name", name), slog.Any("error", err))

		return nil, fmt.Errorf("get %s %q: %w", typ, name, err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "text/markdown",
				Text:     entry.Body,
			},
		},
	}, nil
}

func errorResult(err error) *mcp.CallToolResult {
	return errorResultMsg(err.Error())
}

func errorResultMsg(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
	}
}
