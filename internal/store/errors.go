package store

import "errors"

// ErrNotFound is returned when an entry is not found.
var ErrNotFound = errors.New("not found")

// ErrInvalidFrontmatter is returned when frontmatter is malformed.
var ErrInvalidFrontmatter = errors.New("invalid frontmatter")
