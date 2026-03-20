package models

import "time"

// LogEntry is the canonical representation of a single parsed log line.
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Source    string                 `json:"source"`
	Fields    map[string]interface{} `json:"fields"`
	Raw       string                 `json:"raw"`
}
