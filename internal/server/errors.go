package server

import "errors"

// ErrNameRequired is returned when the name argument is missing.
var ErrNameRequired = errors.New("name is required")
