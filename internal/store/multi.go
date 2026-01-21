package store

import (
	"github.com/monke/grimoire/internal/content"
)

// MultiStore combines multiple stores with fallback behavior.
// It tries each store in order until one succeeds.
type MultiStore struct {
	stores []Store
}

// NewMultiStore creates a store that tries multiple stores in order.
// Typically used to prioritize file store over embedded defaults.
func NewMultiStore(stores ...Store) *MultiStore {
	return &MultiStore{stores: stores}
}

// Get retrieves a single entry, trying each store in order.
func (s *MultiStore) Get(typ content.Type, name string) (*content.Entry, error) {
	var lastErr error

	for _, store := range s.stores {
		entry, err := store.Get(typ, name)
		if err == nil {
			return entry, nil
		}

		lastErr = err
	}

	return nil, lastErr
}

// List returns all entries from all stores, with earlier stores taking precedence.
func (s *MultiStore) List(typ content.Type) ([]*content.Entry, error) {
	seen := make(map[string]bool)

	var result []*content.Entry

	for _, store := range s.stores {
		entries, err := store.List(typ)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !seen[entry.Name] {
				seen[entry.Name] = true
				result = append(result, entry)
			}
		}
	}

	return result, nil
}

// Search finds entries across all stores.
func (s *MultiStore) Search(query string) ([]*content.Entry, error) {
	seen := make(map[string]bool)

	var results []*content.Entry

	for _, store := range s.stores {
		entries, err := store.Search(query)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			key := string(entry.Type) + ":" + entry.Name
			if seen[key] {
				continue
			}

			seen[key] = true

			results = append(results, entry)
		}
	}

	return results, nil
}
