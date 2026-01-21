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

// ErrDuplicate is returned when an entry with the same name already exists.
var ErrDuplicate = errors.New("duplicate entry")

// ErrSourceNotFound is returned when a source path does not exist.
var ErrSourceNotFound = errors.New("source path not found")

// ErrNotDirectory is returned when a source path is not a directory.
var ErrNotDirectory = errors.New("not a directory")

// ErrFilterConflict is returned when both allow and block are configured.
var ErrFilterConflict = errors.New("allow and block cannot both be set")
