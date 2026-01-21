// Package sources provides embedded content for grimoire.
package sources

import "embed"

// FS contains the embedded source files (rules and skills).
//
//go:embed rules skills
var FS embed.FS
