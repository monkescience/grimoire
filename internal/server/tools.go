package server

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/content"
)

// Tool input types.

type emptyInput struct{}

type nameInput struct {
	Name string `json:"name" jsonschema:"Name of the entry"`
}

type searchInput struct {
	Query string `json:"query" jsonschema:"Search query"`
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

// Tool handlers.

func (s *Server) handleListRules(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ emptyInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("listing rules")

	entries := s.store.List(content.TypeRule)

	slog.Debug("rules listed", "count", len(entries))

	return s.entrySummaryResult(entries), nil, nil
}

func (s *Server) handleListPrompts(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ emptyInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("listing prompts")

	entries := s.store.List(content.TypePrompt)

	slog.Debug("prompts listed", "count", len(entries))

	return s.entrySummaryResult(entries), nil, nil
}

func (s *Server) handleListSkills(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ emptyInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("listing skills")

	entries := s.store.List(content.TypeSkill)

	slog.Debug("skills listed", "count", len(entries))

	return s.entrySummaryResult(entries), nil, nil
}

func (s *Server) handleGetRule(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input nameInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("getting rule", "name", input.Name)

	return s.getEntryResult(content.TypeRule, input.Name), nil, nil
}

func (s *Server) handleGetPrompt(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input nameInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("getting prompt", "name", input.Name)

	return s.getEntryResult(content.TypePrompt, input.Name), nil, nil
}

func (s *Server) handleGetSkill(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input nameInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("getting skill", "name", input.Name)

	return s.getEntryResult(content.TypeSkill, input.Name), nil, nil
}

func (s *Server) handleSearch(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input searchInput,
) (*mcp.CallToolResult, any, error) {
	slog.Debug("searching", "query", input.Query)

	entries := s.store.Search(input.Query)

	slog.Debug("search completed", "query", input.Query, "results", len(entries))

	return s.entrySummaryResult(entries), nil, nil
}
