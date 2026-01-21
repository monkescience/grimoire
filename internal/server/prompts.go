package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/content"
)

func (s *Server) registerPrompts() {
	s.mcp.AddPrompt(
		&mcp.Prompt{
			Name:        "apply_rule",
			Description: "Apply a coding rule to guide your work",
			Arguments: []*mcp.PromptArgument{
				{Name: "name", Description: "Name of the rule to apply", Required: true},
			},
		},
		s.handleApplyRulePrompt,
	)

	s.mcp.AddPrompt(
		&mcp.Prompt{
			Name:        "use_skill",
			Description: "Use a skill to perform a specific task",
			Arguments: []*mcp.PromptArgument{
				{Name: "name", Description: "Name of the skill to use", Required: true},
				{Name: "context", Description: "Additional context for the skill", Required: false},
			},
		},
		s.handleUseSkillPrompt,
	)
}

// Prompt handlers.

func (s *Server) handleApplyRulePrompt(
	_ context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	name := req.Params.Arguments["name"]
	if name == "" {
		slog.Warn("apply_rule prompt called without name")

		return nil, ErrNameRequired
	}

	slog.Debug("applying rule prompt", "name", name)

	entry, err := s.store.Get(content.TypeRule, name)
	if err != nil {
		slog.Error("failed to get rule for prompt", "name", name, "error", err)

		return nil, fmt.Errorf("get rule %q: %w", name, err)
	}

	return &mcp.GetPromptResult{
		Description: "Apply rule: " + entry.Title,
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: "Apply the following coding rule to guide your work:\n\n" + entry.Body,
				},
			},
		},
	}, nil
}

func (s *Server) handleUseSkillPrompt(
	_ context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	name := req.Params.Arguments["name"]
	if name == "" {
		slog.Warn("use_skill prompt called without name")

		return nil, ErrNameRequired
	}

	slog.Debug("using skill prompt", "name", name)

	entry, err := s.store.Get(content.TypeSkill, name)
	if err != nil {
		slog.Error("failed to get skill for prompt", "name", name, "error", err)

		return nil, fmt.Errorf("get skill %q: %w", name, err)
	}

	text := entry.Body
	if ctx := req.Params.Arguments["context"]; ctx != "" {
		text = text + "\n\nContext: " + ctx
	}

	return &mcp.GetPromptResult{
		Description: "Use skill: " + entry.Title,
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: text,
				},
			},
		},
	}, nil
}
