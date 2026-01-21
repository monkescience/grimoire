package server

import (
	"context"
	"log/slog"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/content"
)

func (s *Server) registerResources() {
	s.mcp.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        "rule",
			Description: "Get a coding rule by name",
			URITemplate: "grimoire://rules/{name}",
			MIMEType:    "text/markdown",
		},
		s.handleRuleResource,
	)

	s.mcp.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        "prompt",
			Description: "Get a prompt template by name",
			URITemplate: "grimoire://prompts/{name}",
			MIMEType:    "text/markdown",
		},
		s.handlePromptResource,
	)

	s.mcp.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        "skill",
			Description: "Get a skill by name",
			URITemplate: "grimoire://skills/{name}",
			MIMEType:    "text/markdown",
		},
		s.handleSkillResource,
	)
}

// Resource handlers.

func (s *Server) handleRuleResource(
	_ context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://rules/")

	slog.Debug("reading rule resource", "name", name, "uri", req.Params.URI)

	return s.getResourceContents(content.TypeRule, name, req.Params.URI)
}

func (s *Server) handlePromptResource(
	_ context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://prompts/")

	slog.Debug("reading prompt resource", "name", name, "uri", req.Params.URI)

	return s.getResourceContents(content.TypePrompt, name, req.Params.URI)
}

func (s *Server) handleSkillResource(
	_ context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://skills/")

	slog.Debug("reading skill resource", "name", name, "uri", req.Params.URI)

	return s.getResourceContents(content.TypeSkill, name, req.Params.URI)
}
