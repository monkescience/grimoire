// Package sources provides embedded content for grimoire.
package sources

import "embed"

// FS contains the embedded source files (rules, skills, and instructions).
//
//go:embed rules skills instructions
var FS embed.FS
