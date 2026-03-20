package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("logs method path and status", func(t *testing.T) {
		var buf bytes.Buffer
		original := log.Writer()
		log.SetOutput(&buf)
		defer log.SetOutput(original)

		handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		req := httptest.NewRequest(http.MethodDelete, "/api/logs", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		output := buf.String()
		if !strings.Contains(output, "DELETE") {
			t.Fatalf("expected method in log output, got %q", output)
		}
		if !strings.Contains(output, "/api/logs") {
			t.Fatalf("expected path in log output, got %q", output)
		}
		if !strings.Contains(output, "204") {
			t.Fatalf("expected status in log output, got %q", output)
		}
	})
}
