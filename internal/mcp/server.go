package mcp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/monke/grimoire/internal/grimoire"
)

// Server wraps the MCP server with grimoire functionality.
type Server struct {
	mcp   *mcp.Server
	store *grimoire.Store
}

// New creates a new grimoire MCP server.
func New(version string, s *grimoire.Store) *Server {
	srv := &Server{
		store: s,
		mcp: mcp.NewServer(
			&mcp.Implementation{
				Name:    "grimoire",
				Version: version,
			},
			&mcp.ServerOptions{
				Instructions: grimoire.BuildServerInstructions(),
			},
		),
	}

	srv.registerGuidance()
	srv.registerSearch()
	srv.registerResources()

	slog.Debug("server initialized")

	return srv
}

// Run starts the server on stdio.
func (s *Server) Run(ctx context.Context) error {
	slog.Info("starting MCP server on stdio")

	err := s.mcp.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		slog.Error("server stopped with error", slog.Any("error", err))

		return fmt.Errorf("mcp server: %w", err)
	}

	slog.Info("server stopped")

	return nil
}
