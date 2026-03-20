package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("logs exact timestamp method path status and duration format", func(t *testing.T) {
		var buf bytes.Buffer
		original := log.Writer()
		originalFlags := log.Flags()
		log.SetOutput(&buf)
		log.SetFlags(0)
		defer log.SetOutput(original)
		defer log.SetFlags(originalFlags)

		handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		req := httptest.NewRequest(http.MethodDelete, "/api/logs", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		output := buf.String()
		pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z\s+DELETE\s+/api/logs\s+204\s+\d+ms`)
		if !pattern.MatchString(output) {
			t.Fatalf("expected fixed log format, got %q", output)
		}
	})
}
