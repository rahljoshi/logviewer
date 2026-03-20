package parser

import (
	"encoding/json"
	"fmt"
	"io"

	"logviewer/backend/models"
)

// JSONParser parses newline-delimited JSON logs.
type JSONParser struct{}

// ParseFile parses a file using JSON semantics.
func (p JSONParser) ParseFile(r io.Reader, filename string) ([]models.LogEntry, error) {
	return ParseFile(r, filename)
}

// Parse parses a single JSON log line.
func (JSONParser) Parse(line string) (models.LogEntry, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(line), &payload); err != nil {
		return models.LogEntry{}, fmt.Errorf("decode json log line: %w", err)
	}

	entry := models.LogEntry{
		Level:  "INFO",
		Fields: map[string]interface{}{},
	}

	for key, value := range payload {
		switch {
		case hasKey(timestampKeys, key):
			if text, ok := value.(string); ok {
				entry.Timestamp = parseTimestamp(text)
			}
		case hasKey(levelKeys, key):
			if text, ok := value.(string); ok {
				entry.Level = normalizeLevel(text)
			}
		case hasKey(messageKeys, key):
			entry.Message = stringifyValue(value)
		default:
			entry.Fields[key] = value
		}
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = parseTimestamp("")
	}

	return entry, nil
}
