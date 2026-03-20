package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"logviewer/backend/models"
	"logviewer/backend/store"
)

type logsStoreStub struct {
	filtered   []models.LogEntry
	total      int
	stats      map[string]int
	clearCalls int
	filterOpts []store.FilterOpts
}

func (s *logsStoreStub) Add(entries []models.LogEntry) {}

func (s *logsStoreStub) Filter(opts store.FilterOpts) ([]models.LogEntry, int) {
	s.filterOpts = append(s.filterOpts, opts)
	return s.filtered, s.total
}

func (s *logsStoreStub) Stats() map[string]int {
	return s.stats
}

func (s *logsStoreStub) Clear() {
	s.clearCalls++
}

func TestLogsHandler(t *testing.T) {
	t.Run("get logs uses defaults and returns payload", func(t *testing.T) {
		stubStore := &logsStoreStub{
			filtered: []models.LogEntry{{ID: "1", Level: "INFO", Message: "ok"}},
			total:    1,
		}
		handler := LogsHandler{Store: stubStore}
		req := httptest.NewRequest(http.MethodGet, "/api/logs", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, recorder.Code)
		}
		if len(stubStore.filterOpts) != 1 {
			t.Fatalf("expected one filter call, got %d", len(stubStore.filterOpts))
		}
		if stubStore.filterOpts[0].Page != 1 || stubStore.filterOpts[0].PageSize != 100 {
			t.Fatalf("expected default pagination, got %+v", stubStore.filterOpts[0])
		}

		var payload struct {
			Entries  []models.LogEntry `json:"entries"`
			Total    int               `json:"total"`
			Page     int               `json:"page"`
			PageSize int               `json:"pageSize"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Total != 1 || payload.Page != 1 || payload.PageSize != 100 {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})

	t.Run("get logs validates page and page size", func(t *testing.T) {
		handler := LogsHandler{Store: &logsStoreStub{}}

		req := httptest.NewRequest(http.MethodGet, "/api/logs?page=0", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d for invalid page, got %d", http.StatusBadRequest, recorder.Code)
		}

		req = httptest.NewRequest(http.MethodGet, "/api/logs?pageSize=700", nil)
		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d for invalid pageSize, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("delete logs clears store", func(t *testing.T) {
		stubStore := &logsStoreStub{}
		handler := LogsHandler{Store: stubStore}
		req := httptest.NewRequest(http.MethodDelete, "/api/logs", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, recorder.Code)
		}
		if stubStore.clearCalls != 1 {
			t.Fatalf("expected clear once, got %d", stubStore.clearCalls)
		}
	})
}
