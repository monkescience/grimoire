package mcp

import (
	"context"
	"log/slog"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

func (s *Server) registerResources() {
	s.mcp.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        "rule",
			Description: "Get a rule by name",
			URITemplate: "grimoire://rules/{name}",
			MIMEType:    "text/markdown",
		},
		s.handleRuleResource,
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
	ctx context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://rules/")

	slog.DebugContext(ctx, "reading rule resource", slog.String("name", name), slog.String("uri", req.Params.URI))

	return s.getResourceContents(ctx, grimoire.TypeRule, name, req.Params.URI)
}

func (s *Server) handleSkillResource(
	ctx context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://skills/")

	slog.DebugContext(ctx, "reading skill resource", slog.String("name", name), slog.String("uri", req.Params.URI))

	return s.getResourceContents(ctx, grimoire.TypeSkill, name, req.Params.URI)
}
