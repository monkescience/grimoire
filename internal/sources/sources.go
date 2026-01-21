// Package sources provides embedded content for grimoire.
package sources

import "embed"

// FS contains the embedded source files (rules, prompts, skills).
//
//go:embed rules prompts skills
var FS embed.FS
