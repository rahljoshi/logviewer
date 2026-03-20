package parser

import (
	"fmt"
	"io"
	"strings"

	"logviewer/backend/models"
)

// LogfmtParser parses key=value log lines.
type LogfmtParser struct{}

// ParseFile parses a file using logfmt semantics.
func (p LogfmtParser) ParseFile(r io.Reader, filename string) ([]models.LogEntry, error) {
	return ParseFile(r, filename)
}

// Parse parses a single logfmt line.
func (LogfmtParser) Parse(line string) (models.LogEntry, error) {
	fields, err := scanLogfmt(line)
	if err != nil {
		return models.LogEntry{}, err
	}

	entry := models.LogEntry{
		Timestamp: parseTimestamp(""),
		Level:     "INFO",
		Fields:    map[string]interface{}{},
	}

	for key, value := range fields {
		switch {
		case hasKey(timestampKeys, key):
			entry.Timestamp = parseTimestamp(value)
		case hasKey(levelKeys, key):
			entry.Level = normalizeLevel(value)
		case hasKey(messageKeys, key):
			entry.Message = value
		default:
			entry.Fields[key] = value
		}
	}

	return entry, nil
}

func scanLogfmt(line string) (map[string]string, error) {
	fields := make(map[string]string)
	input := strings.TrimSpace(line)
	for len(input) > 0 {
		eq := strings.IndexByte(input, '=')
		if eq <= 0 {
			return nil, fmt.Errorf("invalid logfmt token: %q", input)
		}

		key := input[:eq]
		input = input[eq+1:]

		var value string
		if strings.HasPrefix(input, `"`) {
			input = input[1:]
			var builder strings.Builder
			escaped := false
			closed := false

			for i := 0; i < len(input); i++ {
				ch := input[i]
				if escaped {
					builder.WriteByte(ch)
					escaped = false
					continue
				}
				if ch == '\\' {
					escaped = true
					continue
				}
				if ch == '"' {
					value = builder.String()
					input = input[i+1:]
					closed = true
					break
				}
				builder.WriteByte(ch)
			}

			if !closed {
				return nil, fmt.Errorf("unterminated quoted value for key %q", key)
			}
		} else {
			nextSpace := strings.IndexByte(input, ' ')
			if nextSpace == -1 {
				value = input
				input = ""
			} else {
				value = input[:nextSpace]
				input = input[nextSpace:]
			}
		}

		fields[key] = value
		input = strings.TrimLeft(input, " ")
	}

	return fields, nil
}
