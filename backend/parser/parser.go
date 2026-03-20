package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"logviewer/backend/models"
)

var logfmtPattern = regexp.MustCompile(`\b[\w.-]+=("([^"\\]|\\.)*"|[^\s]+)`)

var timestampKeys = map[string]struct{}{
	"time":       {},
	"ts":         {},
	"timestamp":  {},
	"@timestamp": {},
	"datetime":   {},
}

var levelKeys = map[string]struct{}{
	"level":     {},
	"lvl":       {},
	"severity":  {},
	"loglevel":  {},
	"log_level": {},
}

var messageKeys = map[string]struct{}{
	"msg":     {},
	"message": {},
	"text":    {},
	"body":    {},
	"log":     {},
}

// LogParser parses individual lines and complete files.
type LogParser interface {
	Parse(line string) (models.LogEntry, error)
	ParseFile(r io.Reader, filename string) ([]models.LogEntry, error)
}

// AutoParser detects the input format and delegates to the matching parser.
type AutoParser struct{}

// NewAutoParser constructs an auto-detecting parser.
func NewAutoParser() AutoParser {
	return AutoParser{}
}

// Detect returns the parser kind for a log line.
func Detect(line string) string {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "{") {
		return "json"
	}
	if logfmtPattern.MatchString(line) {
		return "logfmt"
	}
	return "plain"
}

// Parse delegates a single line to the detected parser.
func (AutoParser) Parse(line string) (models.LogEntry, error) {
	return parserForKind(Detect(line)).Parse(line)
}

// ParseFile reads all lines, detects the format once, and parses every entry.
func (AutoParser) ParseFile(r io.Reader, filename string) ([]models.LogEntry, error) {
	return ParseFile(r, filename)
}

// ParseFile reads all lines from r and returns parsed entries.
func ParseFile(r io.Reader, filename string) ([]models.LogEntry, error) {
	_ = filename

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 10<<20), 10<<20)

	entries := make([]models.LogEntry, 0)
	var selected LogParser

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		if selected == nil {
			selected = parserForKind(Detect(line))
		}

		entry, err := selected.Parse(line)
		if err != nil {
			entry = models.LogEntry{
				ID:        uuid.New().String(),
				Timestamp: time.Now(),
				Level:     "UNKNOWN",
				Message:   line,
				Fields:    map[string]interface{}{},
				Raw:       line,
			}
			entries = append(entries, entry)
			continue
		}

		entry.ID = uuid.New().String()
		entry.Raw = line
		if entry.Fields == nil {
			entry.Fields = map[string]interface{}{}
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan log file: %w", err)
	}

	return entries, nil
}

func parserForKind(kind string) LogParser {
	switch kind {
	case "json":
		return JSONParser{}
	case "logfmt":
		return LogfmtParser{}
	default:
		return PlainParser{}
	}
}

func parseTimestamp(value string) time.Time {
	if strings.TrimSpace(value) == "" {
		return time.Now()
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed
		}
	}

	return time.Now()
}

func normalizeLevel(raw string) string {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "INFO":
		return "INFO"
	case "WARN", "WARNING":
		return "WARN"
	case "ERROR", "ERR", "FATAL":
		return "ERROR"
	case "DEBUG", "TRACE":
		return "DEBUG"
	case "":
		return "INFO"
	default:
		return "INFO"
	}
}
