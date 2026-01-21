package grimoire

import "errors"

// ErrNotFound is returned when an entry is not found.
var ErrNotFound = errors.New("not found")

// ErrInvalidFrontmatter is returned when frontmatter is malformed.
var ErrInvalidFrontmatter = errors.New("invalid frontmatter")

// ErrInvalidType is returned when the content type is not recognized.
var ErrInvalidType = errors.New("invalid type")

// ErrNameEmpty is returned when a name is required but empty.
var ErrNameEmpty = errors.New("name is empty")
