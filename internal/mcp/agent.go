package mcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

var errUnexpectedContentType = errors.New("unexpected content type in sampling result")

// agentInput is the input schema for the agent tool.
type agentInput struct {
	Names   []string `json:"names"             jsonschema:"Agent names to execute"`
	Context string   `json:"context,omitempty" jsonschema:"Context provided to agents"`
}

// agentResult holds the output from a single agent execution.
type agentResult struct {
	Name   string
	Output string
	Error  error
}

func (s *Server) registerAgent() {
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "agent",
		Description: "Execute agent prompts via MCP sampling. Requires client sampling support.",
	}, s.handleAgent)
}

func (s *Server) handleAgent(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input agentInput,
) (*mcp.CallToolResult, any, error) {
	slog.DebugContext(ctx, "agent tool called", slog.Any("names", input.Names))

	if len(input.Names) == 0 {
		return errorResultMsg("names parameter is required"), nil, nil
	}

	params := req.Session.InitializeParams()
	if params == nil || params.Capabilities == nil || params.Capabilities.Sampling == nil {
		slog.WarnContext(ctx, "client does not support sampling")

		return s.samplingNotSupported(), nil, nil
	}

	var agents []*grimoire.Entry

	for _, name := range input.Names {
		entry, err := s.store.Get(grimoire.TypeAgent, name)
		if err != nil {
			slog.WarnContext(ctx, "agent not found", slog.String("name", name))

			//nolint:nilerr // return tool error, not Go error
			return errorResultMsg("agent not found: " + name), nil, nil
		}

		agents = append(agents, entry)
	}

	slog.DebugContext(ctx, "executing agents", slog.Int("count", len(agents)))

	var results []agentResult

	for _, agent := range agents {
		slog.DebugContext(ctx, "executing agent", slog.String("name", agent.Name))

		output, err := s.executeSampling(ctx, req.Session, agent, input.Context)
		results = append(results, agentResult{
			Name:   agent.Name,
			Output: output,
			Error:  err,
		})

		if err != nil {
			slog.WarnContext(ctx, "agent execution failed",
				slog.String("name", agent.Name), slog.Any("error", err))
		} else {
			slog.DebugContext(ctx, "agent execution completed", slog.String("name", agent.Name))
		}
	}

	return s.formatAgentResults(results), nil, nil
}

func (s *Server) executeSampling(
	ctx context.Context,
	session *mcp.ServerSession,
	agent *grimoire.Entry,
	context string,
) (string, error) {
	prompt := agent.Body
	if context != "" {
		prompt = prompt + "\n\n## Context\n" + context
	}

	result, err := session.CreateMessage(ctx, &mcp.CreateMessageParams{
		Messages: []*mcp.SamplingMessage{{
			Role:    "user",
			Content: &mcp.TextContent{Text: prompt},
		}},
		SystemPrompt: agent.Description,
	})
	if err != nil {
		return "", fmt.Errorf("sampling failed: %w", err)
	}

	if tc, ok := result.Content.(*mcp.TextContent); ok {
		return tc.Text, nil
	}

	return "", errUnexpectedContentType
}

func (s *Server) samplingNotSupported() *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: `Client does not support MCP sampling.

This feature requires the sampling capability to spawn subagent prompts.
Claude Code support is tracked at: https://github.com/anthropics/claude-code/issues/1785

Alternative: Use guidance(name: "agent-name") to load the prompt and run it manually.`,
			},
		},
	}
}

func (s *Server) formatAgentResults(results []agentResult) *mcp.CallToolResult {
	var b strings.Builder

	for i, r := range results {
		if i > 0 {
			b.WriteString("\n\n---\n\n")
		}

		fmt.Fprintf(&b, "## %s\n\n", r.Name)

		if r.Error != nil {
			fmt.Fprintf(&b, "**Error**: %s\n", r.Error.Error())
		} else {
			b.WriteString(r.Output)
		}
	}

	hasErrors := false

	for _, r := range results {
		if r.Error != nil {
			hasErrors = true

			break
		}
	}

	return &mcp.CallToolResult{
		IsError: hasErrors,
		Content: []mcp.Content{
			&mcp.TextContent{Text: b.String()},
		},
	}
}
