package parser

import (
	"io"
	"regexp"
	"strings"

	"logviewer/backend/models"
)

var plainTimestampPattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}`)
var plainLevelPattern = regexp.MustCompile(`(?i)\b(INFO|WARN|WARNING|ERROR|DEBUG|FATAL|TRACE)\b`)

// PlainParser parses plaintext logs with optional timestamp and level tokens.
type PlainParser struct{}

// ParseFile parses a file using plain text semantics.
func (p PlainParser) ParseFile(r io.Reader, filename string) ([]models.LogEntry, error) {
	return ParseFile(r, filename)
}

// Parse parses a single plain log line.
func (PlainParser) Parse(line string) (models.LogEntry, error) {
	entry := models.LogEntry{
		Timestamp: parseTimestamp(""),
		Level:     "INFO",
		Fields:    map[string]interface{}{},
		Message:   line,
	}

	if timestamp := plainTimestampPattern.FindString(line); timestamp != "" {
		entry.Timestamp = parseTimestamp(timestamp)
	}

	levelLoc := plainLevelPattern.FindStringIndex(line)
	if levelLoc == nil {
		return entry, nil
	}

	entry.Level = normalizeLevel(line[levelLoc[0]:levelLoc[1]])
	entry.Message = strings.TrimSpace(line[levelLoc[1]:])
	if entry.Message == "" {
		entry.Message = line
	}

	return entry, nil
}
