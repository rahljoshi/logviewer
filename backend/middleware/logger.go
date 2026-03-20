package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Round(time.Millisecond)
		if duration < time.Millisecond {
			duration = time.Millisecond
		}

		log.Printf(
			"%s  %-6s  %-20s  %3d  %s",
			start.UTC().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			recorder.status,
			fmt.Sprintf("%dms", duration/time.Millisecond),
		)
	})
}
