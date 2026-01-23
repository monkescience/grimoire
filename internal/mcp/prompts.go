package mcp

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

func (s *Server) registerPrompts() {
	skills := s.store.List(grimoire.TypeSkill)

	for _, skill := range skills {
		s.mcp.AddPrompt(&mcp.Prompt{
			Name:        "grimoire-" + skill.Name,
			Description: skill.Description,
		}, s.makePromptHandler(skill))
	}

	slog.Debug("prompts registered", slog.Int("count", len(skills)))
}

func (s *Server) makePromptHandler(entry *grimoire.Entry) mcp.PromptHandler {
	return func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		slog.DebugContext(ctx, "prompt requested", slog.String("name", entry.Name))

		return &mcp.GetPromptResult{
			Description: entry.Description,
			Messages: []*mcp.PromptMessage{
				{
					Role:    "user",
					Content: &mcp.TextContent{Text: entry.Body},
				},
			},
		}, nil
	}
}
