// Package sources provides embedded content for grimoire.
package sources

import "embed"

// FS contains the embedded source files (rules, skills, instructions, and agents).
//
//go:embed rules skills instructions agents
var FS embed.FS
