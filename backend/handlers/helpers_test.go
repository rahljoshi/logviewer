package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	t.Run("writes json body and content type", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		respondJSON(recorder, http.StatusCreated, map[string]string{"ok": "true"})

		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
		}
		if got := recorder.Header().Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected application/json content type, got %q", got)
		}

		var payload map[string]string
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("expected valid json body: %v", err)
		}
		if payload["ok"] != "true" {
			t.Fatalf("expected json payload, got %#v", payload)
		}
	})
}

func TestRespondError(t *testing.T) {
	t.Run("writes standard error envelope", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		respondError(recorder, http.StatusBadRequest, "missing file")

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
		}

		var payload map[string]string
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("expected valid json body: %v", err)
		}
		if payload["error"] != "missing file" {
			t.Fatalf("expected error envelope, got %#v", payload)
		}
	})
}
