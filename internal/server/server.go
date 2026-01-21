package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/content"
	"github.com/monke/grimoire/internal/store"
)

// ErrNameRequired is returned when the name argument is missing.
var ErrNameRequired = errors.New("name is required")

// Server wraps the MCP server with grimoire functionality.
type Server struct {
	mcp   *mcp.Server
	store store.Store
}

// New creates a new grimoire MCP server.
func New(version string, s store.Store) *Server {
	srv := &Server{
		store: s,
		mcp: mcp.NewServer(
			&mcp.Implementation{
				Name:    "grimoire",
				Version: version,
			},
			&mcp.ServerOptions{
				Instructions: "Grimoire is a knowledge server providing coding rules, prompts, and skills.",
			},
		),
	}

	srv.registerTools()
	srv.registerResources()
	srv.registerPrompts()

	return srv
}

// Run starts the server on stdio.
func (s *Server) Run(ctx context.Context) error {
	err := s.mcp.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		return fmt.Errorf("mcp server: %w", err)
	}

	return nil
}

func (s *Server) registerTools() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_rules",
		Description: "List all available coding rules",
	}, s.handleListRules)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_prompts",
		Description: "List all available prompts",
	}, s.handleListPrompts)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_skills",
		Description: "List all available skills",
	}, s.handleListSkills)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_rule",
		Description: "Get a specific coding rule by name",
	}, s.handleGetRule)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_prompt",
		Description: "Get a specific prompt by name",
	}, s.handleGetPrompt)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_skill",
		Description: "Get a specific skill by name",
	}, s.handleGetSkill)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "search",
		Description: "Search across all rules, prompts, and skills",
	}, s.handleSearch)
}

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

// Tool input types.

type emptyInput struct{}

type nameInput struct {
	Name string `json:"name" jsonschema:"Name of the entry"`
}

type searchInput struct {
	Query string `json:"query" jsonschema:"Search query"`
}

// Tool handlers.

func (s *Server) handleListRules(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ emptyInput,
) (*mcp.CallToolResult, any, error) {
	entries, err := s.store.List(content.TypeRule)
	if err != nil {
		return errorResult(err), nil, nil
	}

	return s.entrySummaryResult(entries), nil, nil
}

func (s *Server) handleListPrompts(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ emptyInput,
) (*mcp.CallToolResult, any, error) {
	entries, err := s.store.List(content.TypePrompt)
	if err != nil {
		return errorResult(err), nil, nil
	}

	return s.entrySummaryResult(entries), nil, nil
}

func (s *Server) handleListSkills(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ emptyInput,
) (*mcp.CallToolResult, any, error) {
	entries, err := s.store.List(content.TypeSkill)
	if err != nil {
		return errorResult(err), nil, nil
	}

	return s.entrySummaryResult(entries), nil, nil
}

func (s *Server) handleGetRule(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input nameInput,
) (*mcp.CallToolResult, any, error) {
	return s.getEntryResult(content.TypeRule, input.Name), nil, nil
}

func (s *Server) handleGetPrompt(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input nameInput,
) (*mcp.CallToolResult, any, error) {
	return s.getEntryResult(content.TypePrompt, input.Name), nil, nil
}

func (s *Server) handleGetSkill(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input nameInput,
) (*mcp.CallToolResult, any, error) {
	return s.getEntryResult(content.TypeSkill, input.Name), nil, nil
}

func (s *Server) handleSearch(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input searchInput,
) (*mcp.CallToolResult, any, error) {
	entries, err := s.store.Search(input.Query)
	if err != nil {
		return errorResult(err), nil, nil
	}

	return s.entrySummaryResult(entries), nil, nil
}

// Resource handlers.

func (s *Server) handleRuleResource(
	_ context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://rules/")

	return s.getResourceContents(content.TypeRule, name, req.Params.URI)
}

func (s *Server) handlePromptResource(
	_ context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://prompts/")

	return s.getResourceContents(content.TypePrompt, name, req.Params.URI)
}

func (s *Server) handleSkillResource(
	_ context.Context,
	req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	name := strings.TrimPrefix(req.Params.URI, "grimoire://skills/")

	return s.getResourceContents(content.TypeSkill, name, req.Params.URI)
}

// Prompt handlers.

func (s *Server) handleApplyRulePrompt(
	_ context.Context,
	req *mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	name := req.Params.Arguments["name"]
	if name == "" {
		return nil, ErrNameRequired
	}

	entry, err := s.store.Get(content.TypeRule, name)
	if err != nil {
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
		return nil, ErrNameRequired
	}

	entry, err := s.store.Get(content.TypeSkill, name)
	if err != nil {
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

// Helper types and functions.

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
