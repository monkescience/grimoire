package server

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/content"
)

// entrySummary is a lightweight representation of an entry for list responses.
type entrySummary struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (s *Server) entrySummaryResult(entries []*content.Entry) *mcp.CallToolResult {
	summaries := make([]entrySummary, len(entries))
	for i, e := range entries {
		summaries[i] = entrySummary{
			Name:        e.Name,
			Type:        string(e.Type),
			Title:       e.Title,
			Description: e.Description,
			Tags:        e.Tags,
		}
	}

	data, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		slog.Error("failed to marshal entry summaries", "error", err)

		return errorResult(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}
}

func (s *Server) getEntryResult(typ content.Type, name string) *mcp.CallToolResult {
	entry, err := s.store.Get(typ, name)
	if err != nil {
		slog.Warn("failed to get entry", "type", typ, "name", name, "error", err)

		return errorResult(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: entry.Body},
		},
	}
}

func (s *Server) getResourceContents(
	typ content.Type,
	name, uri string,
) (*mcp.ReadResourceResult, error) {
	entry, err := s.store.Get(typ, name)
	if err != nil {
		slog.Warn("failed to get resource contents", "type", typ, "name", name, "error", err)

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
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: err.Error()},
		},
	}
}
