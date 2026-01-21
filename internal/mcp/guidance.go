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
	Name   string   `json:"name,omitempty"   jsonschema:"Name of the guidance to load (exact match)"`
	Topics []string `json:"topics,omitempty" jsonschema:"Topics to find matching rules (e.g., [go, errors])"`
	Files  []string `json:"files,omitempty"  jsonschema:"File paths to find matching rules by globs (e.g., [main.go])"`
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
	// Handle exact name lookup (existing behavior)
	if input.Name != "" {
		return s.handleGuidanceByName(input.Name)
	}

	// Handle topic-based lookup
	if len(input.Topics) > 0 {
		return s.handleGuidanceByTopics(input.Topics)
	}

	// Handle file-based lookup
	if len(input.Files) > 0 {
		return s.handleGuidanceByFiles(input.Files)
	}

	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: "provide one of: name, topics, or files"},
		},
	}, nil, nil
}

func (s *Server) handleGuidanceByName(name string) (*mcp.CallToolResult, any, error) {
	slog.Debug("loading guidance by name", "name", name)

	// Try each type until we find a match
	for _, typ := range []grimoire.Type{grimoire.TypeSkill, grimoire.TypeRule} {
		entry, err := s.store.Get(typ, name)
		if err == nil {
			slog.Debug("guidance loaded", "name", name, "type", typ)

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: entry.Body},
				},
			}, nil, nil
		}
	}

	slog.Warn("guidance not found", "name", name)

	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("guidance %q not found", name)},
		},
	}, nil, nil
}

func (s *Server) handleGuidanceByTopics(topics []string) (*mcp.CallToolResult, any, error) {
	slog.Debug("loading guidance by topics", "topics", topics)

	entries := s.store.FindByTopics(topics)
	if len(entries) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("no rules found for topics: %v", topics)},
			},
		}, nil, nil
	}

	slog.Debug("guidance loaded by topics", "topics", topics, "count", len(entries))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatEntries(entries)},
		},
	}, nil, nil
}

func (s *Server) handleGuidanceByFiles(files []string) (*mcp.CallToolResult, any, error) {
	slog.Debug("loading guidance by files", "files", files)

	entries := s.store.FindByGlobs(files)
	if len(entries) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("no rules found for files: %v", files)},
			},
		}, nil, nil
	}

	slog.Debug("guidance loaded by files", "files", files, "count", len(entries))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: formatEntries(entries)},
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
