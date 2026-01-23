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
			Name:        skill.Name,
			Description: skill.Description,
			Arguments:   convertArguments(skill.Arguments),
		}, s.makePromptHandler(skill))
	}

	slog.Debug("prompts registered", slog.Int("count", len(skills)))
}

func convertArguments(args []grimoire.Argument) []*mcp.PromptArgument {
	if len(args) == 0 {
		return nil
	}

	result := make([]*mcp.PromptArgument, len(args))
	for i, arg := range args {
		result[i] = &mcp.PromptArgument{
			Name:        arg.Name,
			Description: arg.Description,
			Required:    arg.Required,
		}
	}

	return result
}

func (s *Server) makePromptHandler(entry *grimoire.Entry) mcp.PromptHandler {
	return func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		slog.DebugContext(ctx, "prompt requested",
			slog.String("name", entry.Name),
			slog.Any("arguments", req.Params.Arguments))

		body := entry.RenderBody(req.Params.Arguments)

		return &mcp.GetPromptResult{
			Description: entry.Description,
			Messages: []*mcp.PromptMessage{
				{
					Role:    "user",
					Content: &mcp.TextContent{Text: body},
				},
			},
		}, nil
	}
}
