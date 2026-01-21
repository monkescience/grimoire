package store

import (
	"github.com/monke/grimoire/internal/content"
)

// Store provides access to grimoire content.
type Store interface {
	// Get retrieves a single entry by type and name.
	Get(typ content.Type, name string) (*content.Entry, error)

	// List returns all entries of a given type.
	List(typ content.Type) ([]*content.Entry, error)

	// Search finds entries matching the query across all types.
	Search(query string) ([]*content.Entry, error)
}
