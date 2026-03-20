package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatsHandler(t *testing.T) {
	t.Run("returns total and level counts", func(t *testing.T) {
		stubStore := &logsStoreStub{
			stats: map[string]int{
				"INFO":  2,
				"ERROR": 1,
			},
		}
		handler := StatsHandler{Store: stubStore}
		req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, recorder.Code)
		}

		var payload struct {
			Total  int            `json:"total"`
			Levels map[string]int `json:"levels"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Total != 3 {
			t.Fatalf("expected total 3, got %d", payload.Total)
		}
		if payload.Levels["INFO"] != 2 || payload.Levels["ERROR"] != 1 {
			t.Fatalf("unexpected levels payload: %#v", payload.Levels)
		}
	})
}
