package grimoire

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for loading grimoire content.
type Config struct {
	// Sources configures where to load content from.
	Sources SourcesConfig `yaml:"sources"`

	// Rules configures filtering for rules.
	Rules FilterConfig `yaml:"rules"`

	// Skills configures filtering for skills.
	Skills FilterConfig `yaml:"skills"`
}

// SourcesConfig configures content sources.
type SourcesConfig struct {
	// Builtin enables loading embedded content. Default: true.
	Builtin *bool `yaml:"builtin"`

	// Paths lists external directories to load content from.
	// Paths are checked in order; duplicates cause an error.
	Paths []string `yaml:"paths"`
}

// FilterConfig configures allow/block filtering for a content type.
type FilterConfig struct {
	// Allow lists names to allow. If non-empty, only these are loaded.
	Allow []string `yaml:"allow"`

	// Block lists names to block. Ignored if Allow is non-empty.
	Block []string `yaml:"block"`
}

// BuiltinEnabled returns whether builtin content should be loaded.
func (c *Config) BuiltinEnabled() bool {
	if c.Sources.Builtin == nil {
		return true
	}

	return *c.Sources.Builtin
}

// DefaultConfig returns the default configuration (builtin only).
func DefaultConfig() *Config {
	return &Config{}
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	path = ExpandHome(path)

	data, err := os.ReadFile(path) //nolint:gosec // User-provided config path is intentional
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Expand ~ in all paths
	for i, p := range cfg.Sources.Paths {
		cfg.Sources.Paths[i] = ExpandHome(p)
	}

	return &cfg, nil
}

// ExpandHome expands ~ to the user's home directory.
func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, path[1:])
}

// IsAllowed checks if a name is allowed by the filter config.
func (f *FilterConfig) IsAllowed(name string) bool {
	// If allow list is set, only those names are allowed
	if len(f.Allow) > 0 {
		return slices.Contains(f.Allow, name)
	}

	// Otherwise, check block list
	return !slices.Contains(f.Block, name)
}

// FilterForType returns the filter config for a given content type.
func (c *Config) FilterForType(typ Type) *FilterConfig {
	switch typ {
	case TypeRule:
		return &c.Rules
	case TypeSkill:
		return &c.Skills
	default:
		return &FilterConfig{}
	}
}
