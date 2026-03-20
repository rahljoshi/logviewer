package store

import (
	"testing"
	"time"

	"logviewer/backend/models"
)

func TestMemStoreFilter(t *testing.T) {
	t.Run("filter by level returns only matching entries", func(t *testing.T) {
		store := NewMemStore()
		store.Add(sampleEntries())

		got, total := store.Filter(FilterOpts{
			Level:    "ERROR",
			Page:     1,
			PageSize: 10,
		})

		if total != 1 {
			t.Fatalf("expected total 1, got %d", total)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(got))
		}
		if got[0].Level != "ERROR" {
			t.Fatalf("expected ERROR level, got %s", got[0].Level)
		}
	})

	t.Run("filter by query is case-insensitive", func(t *testing.T) {
		store := NewMemStore()
		store.Add(sampleEntries())

		got, total := store.Filter(FilterOpts{
			Query:    "TIMED OUT",
			Page:     1,
			PageSize: 10,
		})

		if total != 1 {
			t.Fatalf("expected total 1, got %d", total)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(got))
		}
		if got[0].Message != "request timed out" {
			t.Fatalf("expected timeout message, got %q", got[0].Message)
		}
	})

	t.Run("pagination returns correct window and total", func(t *testing.T) {
		store := NewMemStore()
		store.Add(sampleEntries())

		got, total := store.Filter(FilterOpts{
			Page:     2,
			PageSize: 2,
		})

		if total != 4 {
			t.Fatalf("expected total 4, got %d", total)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(got))
		}
		if got[0].ID != "3" || got[1].ID != "4" {
			t.Fatalf("expected entries 3 and 4, got %q and %q", got[0].ID, got[1].ID)
		}
	})
}

func TestMemStoreStats(t *testing.T) {
	t.Run("stats counts each level correctly", func(t *testing.T) {
		store := NewMemStore()
		store.Add(sampleEntries())

		got := store.Stats()

		expected := map[string]int{
			"INFO":  1,
			"WARN":  1,
			"ERROR": 1,
			"DEBUG": 1,
		}

		for level, want := range expected {
			if got[level] != want {
				t.Fatalf("expected %s count %d, got %d", level, want, got[level])
			}
		}
	})
}

func TestMemStoreClear(t *testing.T) {
	t.Run("clear empties the store", func(t *testing.T) {
		store := NewMemStore()
		store.Add(sampleEntries())

		store.Clear()

		got, total := store.Filter(FilterOpts{
			Page:     1,
			PageSize: 10,
		})

		if total != 0 {
			t.Fatalf("expected total 0, got %d", total)
		}
		if len(got) != 0 {
			t.Fatalf("expected 0 entries, got %d", len(got))
		}
	})
}

func sampleEntries() []models.LogEntry {
	base := time.Date(2026, 3, 21, 1, 0, 0, 0, time.UTC)

	return []models.LogEntry{
		{ID: "1", Timestamp: base, Level: "INFO", Message: "user login successful", Source: "app.log"},
		{ID: "2", Timestamp: base.Add(time.Second), Level: "WARN", Message: "database query slow", Source: "app.log"},
		{ID: "3", Timestamp: base.Add(2 * time.Second), Level: "ERROR", Message: "request timed out", Source: "app.log"},
		{ID: "4", Timestamp: base.Add(3 * time.Second), Level: "DEBUG", Message: "cache miss for key alpha", Source: "app.log"},
	}
}
