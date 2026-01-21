package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/monke/grimoire/internal/server"
	"github.com/monke/grimoire/internal/sources"
	"github.com/monke/grimoire/internal/store"
)

var version = "dev"

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	contentDir := flag.String("content", "", "Path to content directory (overrides embedded content)")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersion {
		_, _ = fmt.Fprintln(os.Stdout, version)

		return nil
	}

	// Create store based on configuration
	var s store.Store

	embedStore := store.NewEmbedStore(sources.FS, "")

	if *contentDir != "" {
		// Use file store with embedded fallback
		fileStore := store.NewFileStore(*contentDir)
		s = store.NewMultiStore(fileStore, embedStore)
	} else {
		// Use embedded content only
		s = embedStore
	}

	// Create and run server
	srv := server.New(version, s)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	err := srv.Run(ctx)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	return nil
}
