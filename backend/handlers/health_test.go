package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	t.Run("returns ok true", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		recorder := httptest.NewRecorder()

		HealthHandler{}.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, recorder.Code)
		}

		var payload map[string]bool
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if !payload["ok"] {
			t.Fatalf("expected ok response, got %#v", payload)
		}
	})
}
