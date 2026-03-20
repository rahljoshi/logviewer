package store

import (
	"strings"
	"sync"

	"logviewer/backend/models"
)

// Store describes the in-memory log entry storage contract.
type Store interface {
	Add(entries []models.LogEntry)
	Filter(opts FilterOpts) ([]models.LogEntry, int)
	Stats() map[string]int
	Clear()
}

// FilterOpts controls level, message, and paging filters.
type FilterOpts struct {
	Level    string
	Query    string
	Page     int
	PageSize int
}

// MemStore is a thread-safe in-memory implementation of Store.
type MemStore struct {
	mu      sync.RWMutex
	entries []models.LogEntry
}

// NewMemStore creates an empty MemStore.
func NewMemStore() *MemStore {
	return &MemStore{
		entries: make([]models.LogEntry, 0),
	}
}

// Add appends entries to the store.
func (s *MemStore) Add(entries []models.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = append(s.entries, entries...)
}

// Filter returns the current page of matching entries and the total match count.
func (s *MemStore) Filter(opts FilterOpts) ([]models.LogEntry, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	page := opts.Page
	if page < 1 {
		page = 1
	}

	pageSize := opts.PageSize
	if pageSize < 1 {
		pageSize = len(s.entries)
	}

	query := strings.ToLower(opts.Query)
	filtered := make([]models.LogEntry, 0, len(s.entries))

	for _, entry := range s.entries {
		if opts.Level != "" && entry.Level != opts.Level {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(entry.Message), query) {
			continue
		}
		filtered = append(filtered, entry)
	}

	total := len(filtered)
	start := (page - 1) * pageSize
	if start >= total {
		return []models.LogEntry{}, total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	pageEntries := make([]models.LogEntry, end-start)
	copy(pageEntries, filtered[start:end])

	return pageEntries, total
}

// Stats returns counts by level for all entries in the store.
func (s *MemStore) Stats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]int)
	for _, entry := range s.entries {
		stats[entry.Level]++
	}

	return stats
}

// Clear removes all entries from the store.
func (s *MemStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = s.entries[:0]
}
