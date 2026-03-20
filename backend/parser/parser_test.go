package parser

import (
	"strings"
	"testing"
	"time"
)

func TestDetect(t *testing.T) {
	t.Run("detects json format", func(t *testing.T) {
		got := Detect(`{"level":"info","message":"hello"}`)
		if got != "json" {
			t.Fatalf("expected json, got %q", got)
		}
	})

	t.Run("detects logfmt format", func(t *testing.T) {
		got := Detect(`level=info msg="hello world"`)
		if got != "logfmt" {
			t.Fatalf("expected logfmt, got %q", got)
		}
	})

	t.Run("detects plain format", func(t *testing.T) {
		got := Detect(`2026-03-21 01:00:00 INFO service started`)
		if got != "plain" {
			t.Fatalf("expected plain, got %q", got)
		}
	})
}

func TestJSONParser(t *testing.T) {
	parser := JSONParser{}

	t.Run("valid line maps canonical fields", func(t *testing.T) {
		entry, err := parser.Parse(`{"timestamp":"2026-03-21T01:02:03Z","level":"warning","message":"database query slow","service":"api","request_id":"abc-123"}`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Level != "WARN" {
			t.Fatalf("expected WARN, got %q", entry.Level)
		}
		if entry.Message != "database query slow" {
			t.Fatalf("expected mapped message, got %q", entry.Message)
		}
		if entry.Timestamp.Format(time.RFC3339) != "2026-03-21T01:02:03Z" {
			t.Fatalf("unexpected timestamp %s", entry.Timestamp.Format(time.RFC3339))
		}
		if entry.Source != "" {
			t.Fatalf("expected parser to leave source empty, got %q", entry.Source)
		}
		if got := entry.Fields["service"]; got != "api" {
			t.Fatalf("expected service field preserved, got %#v", got)
		}
		if got := entry.Fields["request_id"]; got != "abc-123" {
			t.Fatalf("expected request_id field preserved, got %#v", got)
		}
	})

	t.Run("missing timestamp uses current time", func(t *testing.T) {
		before := time.Now().Add(-time.Second)
		entry, err := parser.Parse(`{"level":"info","message":"user login successful"}`)
		after := time.Now().Add(time.Second)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
			t.Fatalf("expected timestamp to be set near now, got %s", entry.Timestamp)
		}
	})

	t.Run("missing level defaults to info", func(t *testing.T) {
		entry, err := parser.Parse(`{"message":"cache warm complete"}`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Level != "INFO" {
			t.Fatalf("expected INFO, got %q", entry.Level)
		}
	})

	t.Run("unknown fields are preserved in fields map", func(t *testing.T) {
		entry, err := parser.Parse(`{"msg":"request completed","region":"ap-south-1","status":200}`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got := entry.Fields["region"]; got != "ap-south-1" {
			t.Fatalf("expected region field preserved, got %#v", got)
		}
		if got := entry.Fields["status"]; got != float64(200) {
			t.Fatalf("expected status field preserved, got %#v", got)
		}
	})
}

func TestLogfmtParser(t *testing.T) {
	parser := LogfmtParser{}

	t.Run("simple line parses key value pairs", func(t *testing.T) {
		entry, err := parser.Parse(`time=2026-03-21T01:02:03Z level=error msg="request failed" code=500`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Level != "ERROR" {
			t.Fatalf("expected ERROR, got %q", entry.Level)
		}
		if entry.Message != "request failed" {
			t.Fatalf("expected mapped message, got %q", entry.Message)
		}
		if got := entry.Fields["code"]; got != "500" {
			t.Fatalf("expected code field preserved, got %#v", got)
		}
	})

	t.Run("quoted values are unwrapped", func(t *testing.T) {
		entry, err := parser.Parse(`ts="2026-03-21 01:02:03" lvl=warn msg="cache miss for key alpha" component="edge proxy"`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Level != "WARN" {
			t.Fatalf("expected WARN, got %q", entry.Level)
		}
		if got := entry.Fields["component"]; got != "edge proxy" {
			t.Fatalf("expected component field preserved, got %#v", got)
		}
	})

	t.Run("missing level defaults to info", func(t *testing.T) {
		entry, err := parser.Parse(`msg="request queued" trace_id=xyz`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Level != "INFO" {
			t.Fatalf("expected INFO, got %q", entry.Level)
		}
	})
}

func TestPlainParser(t *testing.T) {
	parser := PlainParser{}

	t.Run("iso timestamp and level are extracted", func(t *testing.T) {
		entry, err := parser.Parse(`2026-03-21T01:02:03 ERROR request timed out after 3s`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Level != "ERROR" {
			t.Fatalf("expected ERROR, got %q", entry.Level)
		}
		if entry.Message != "request timed out after 3s" {
			t.Fatalf("expected trimmed message, got %q", entry.Message)
		}
	})

	t.Run("missing timestamp uses current time", func(t *testing.T) {
		before := time.Now().Add(-time.Second)
		entry, err := parser.Parse(`WARN cache pressure rising`)
		after := time.Now().Add(time.Second)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
			t.Fatalf("expected current timestamp, got %s", entry.Timestamp)
		}
	})

	t.Run("missing level keeps whole line as message", func(t *testing.T) {
		entry, err := parser.Parse(`2026-03-21 01:02:03 background worker heartbeat`)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entry.Level != "INFO" {
			t.Fatalf("expected INFO default, got %q", entry.Level)
		}
		if entry.Message != "2026-03-21 01:02:03 background worker heartbeat" {
			t.Fatalf("expected whole line as message, got %q", entry.Message)
		}
	})
}

func TestParseFile(t *testing.T) {
	t.Run("multi line file returns all entries for detected format", func(t *testing.T) {
		input := strings.Join([]string{
			`{"timestamp":"2026-03-21T01:00:00Z","level":"info","message":"user login successful"}`,
			`{"timestamp":"2026-03-21T01:00:02Z","level":"error","message":"request timed out"}`,
		}, "\n")

		entries, err := ParseFile(strings.NewReader(input), "app.log")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}
		for _, entry := range entries {
			if entry.ID == "" {
				t.Fatal("expected entry id to be populated")
			}
			if entry.Source != "" {
				t.Fatalf("expected source to remain empty, got %q", entry.Source)
			}
		}
	})

	t.Run("empty lines are skipped", func(t *testing.T) {
		input := "\n\nlevel=info msg=\"ok\"\n  \nlevel=warn msg=\"slow\"\n"

		entries, err := ParseFile(strings.NewReader(input), "app.log")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}
	})

	t.Run("bad line produces unknown entry", func(t *testing.T) {
		input := strings.Join([]string{
			`level=info msg="ok"`,
			`msg="unterminated`,
		}, "\n")

		entries, err := ParseFile(strings.NewReader(input), "app.log")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}
		if entries[1].Level != "UNKNOWN" {
			t.Fatalf("expected UNKNOWN fallback, got %q", entries[1].Level)
		}
		if entries[1].Message != `msg="unterminated` {
			t.Fatalf("expected raw line as message, got %q", entries[1].Message)
		}
		if entries[1].Raw != `msg="unterminated` {
			t.Fatalf("expected raw line preserved, got %q", entries[1].Raw)
		}
	})
}
