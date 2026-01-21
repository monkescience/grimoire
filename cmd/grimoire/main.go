package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/monke/grimoire/internal/grimoire"
	"github.com/monke/grimoire/internal/mcp"
	"github.com/monke/grimoire/internal/sources"
)

var version = "dev"

func main() {
	err := run()
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	showVersion := flag.Bool("version", false, "Show version")
	verbose := flag.Bool("verbose", false, "Enable verbose logging (debug level)")

	flag.Parse()

	if *showVersion {
		_, _ = fmt.Fprintln(os.Stdout, version)

		return nil
	}

	// Configure logger based on verbose flag
	configureLogger(*verbose)

	slog.Info("starting grimoire", "version", version)

	// Create store from embedded sources
	s, err := grimoire.New(sources.FS)
	if err != nil {
		return fmt.Errorf("loading sources: %w", err)
	}

	slog.Debug("store initialized")

	// Create and run server
	srv := mcp.New(version, s)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	err = srv.Run(ctx)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	return nil
}

func configureLogger(verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	// Write logs to stderr to keep stdout clean for MCP protocol
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))
}
