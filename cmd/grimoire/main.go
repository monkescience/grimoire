package main

import (
	"context"
	"errors"
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

// errConfigConflict is returned when --config is combined with other flags.
var errConfigConflict = errors.New("--config cannot be combined with --source, --no-builtin, or filter flags")

// stringSlice implements flag.Value for repeated string flags.
type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)

	return nil
}

// flags holds all parsed command-line flags.
type flags struct {
	showVersion bool
	verbose     bool
	configFile  string
	sourcePaths stringSlice
	noBuiltin   bool
	allowRules  stringSlice
	blockRules  stringSlice
	allowSkills stringSlice
	blockSkills stringSlice
}

func main() {
	err := run()
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	f := parseFlags()

	if f.showVersion {
		_, _ = fmt.Fprintln(os.Stdout, version)

		return nil
	}

	configureLogger(f.verbose)

	cfg, err := buildConfig(f)
	if err != nil {
		return err
	}

	slog.Info("starting grimoire", "version", version)

	return runServer(cfg)
}

func parseFlags() *flags {
	f := &flags{}

	flag.BoolVar(&f.showVersion, "version", false, "Show version")
	flag.BoolVar(&f.verbose, "verbose", false, "Enable verbose logging (debug level)")
	flag.StringVar(&f.configFile, "config", "", "Load configuration from YAML file")
	flag.Var(&f.sourcePaths, "source", "External source directory (can be repeated)")
	flag.BoolVar(&f.noBuiltin, "no-builtin", false, "Disable embedded/builtin content")
	flag.Var(&f.allowRules, "allow-rule", "Only load these rules (can be repeated)")
	flag.Var(&f.blockRules, "block-rule", "Block these rules (can be repeated)")
	flag.Var(&f.allowSkills, "allow-skill", "Only load these skills (can be repeated)")
	flag.Var(&f.blockSkills, "block-skill", "Block these skills (can be repeated)")

	flag.Parse()

	return f
}

func runServer(cfg *grimoire.Config) error {
	store, err := grimoire.New(cfg, sources.FS)
	if err != nil {
		return fmt.Errorf("loading sources: %w", err)
	}

	slog.Debug("store initialized")

	srv := mcp.New(version, store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

// buildConfig creates a Config from either a config file or CLI flags.
// Config file and CLI flags are mutually exclusive.
func buildConfig(f *flags) (*grimoire.Config, error) {
	hasCLIFlags := len(f.sourcePaths) > 0 || f.noBuiltin ||
		len(f.allowRules) > 0 || len(f.blockRules) > 0 ||
		len(f.allowSkills) > 0 || len(f.blockSkills) > 0

	if f.configFile != "" && hasCLIFlags {
		return nil, errConfigConflict
	}

	if f.configFile != "" {
		slog.Debug("loading config from file", "path", f.configFile)

		cfg, err := grimoire.LoadConfig(f.configFile)
		if err != nil {
			return nil, fmt.Errorf("loading config: %w", err)
		}

		return cfg, nil
	}

	cfg := grimoire.DefaultConfig()

	if f.noBuiltin {
		builtin := false
		cfg.Sources.Builtin = &builtin
	}

	cfg.Sources.Paths = expandPaths(f.sourcePaths)
	cfg.Rules.Allow = f.allowRules
	cfg.Rules.Block = f.blockRules
	cfg.Skills.Allow = f.allowSkills
	cfg.Skills.Block = f.blockSkills

	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// expandPaths expands ~ in all paths.
func expandPaths(paths []string) []string {
	result := make([]string, len(paths))
	for i, p := range paths {
		result[i] = grimoire.ExpandHome(p)
	}

	return result
}

func configureLogger(verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))
}
