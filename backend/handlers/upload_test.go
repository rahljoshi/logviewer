package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"logviewer/backend/models"
	"logviewer/backend/store"
)

type parserStub struct {
	entries []models.LogEntry
	err     error
}

func (s parserStub) Parse(line string) (models.LogEntry, error) {
	return models.LogEntry{}, nil
}

func (s parserStub) ParseFile(_ io.Reader, _ string) ([]models.LogEntry, error) {
	return s.entries, s.err
}

type storeStub struct {
	added      []models.LogEntry
	stats      map[string]int
	clearCalls int
}

func (s *storeStub) Add(entries []models.LogEntry) {
	s.added = append([]models.LogEntry(nil), entries...)
}

func (s *storeStub) Filter(opts store.FilterOpts) ([]models.LogEntry, int) {
	return nil, 0
}

func (s *storeStub) Stats() map[string]int {
	return s.stats
}

func (s *storeStub) Clear() {
	s.clearCalls++
}

func TestUploadHandler(t *testing.T) {
	t.Run("missing file returns bad request", func(t *testing.T) {
		handler := UploadHandler{
			Parser: parserStub{},
			Store:  &storeStub{},
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		if err := writer.Close(); err != nil {
			t.Fatalf("close multipart writer: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("parser failure returns unprocessable entity", func(t *testing.T) {
		handler := UploadHandler{
			Parser: parserStub{err: errors.New("bad log file")},
			Store:  &storeStub{},
		}
		req := newUploadRequest(t, "app.log", "payload")
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected %d, got %d", http.StatusUnprocessableEntity, recorder.Code)
		}
	})

	t.Run("successful upload clears previous entries and sets source", func(t *testing.T) {
		stubStore := &storeStub{
			stats: map[string]int{"INFO": 1, "ERROR": 1},
		}
		handler := UploadHandler{
			Parser: parserStub{entries: []models.LogEntry{
				{ID: "1", Timestamp: time.Now(), Level: "INFO", Message: "ok", Fields: map[string]interface{}{}, Raw: "ok"},
				{ID: "2", Timestamp: time.Now(), Level: "ERROR", Message: "bad", Fields: map[string]interface{}{}, Raw: "bad"},
			}},
			Store: stubStore,
		}
		req := newUploadRequest(t, "server.log", "payload")
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, recorder.Code)
		}
		if stubStore.clearCalls != 1 {
			t.Fatalf("expected clear to be called once, got %d", stubStore.clearCalls)
		}
		if len(stubStore.added) != 2 {
			t.Fatalf("expected 2 entries added, got %d", len(stubStore.added))
		}
		for _, entry := range stubStore.added {
			if entry.Source != "server.log" {
				t.Fatalf("expected source to be set from filename, got %q", entry.Source)
			}
		}

		var payload struct {
			Count int            `json:"count"`
			Stats map[string]int `json:"stats"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Count != 2 {
			t.Fatalf("expected count 2, got %d", payload.Count)
		}
		if payload.Stats["ERROR"] != 1 {
			t.Fatalf("expected stats payload, got %#v", payload.Stats)
		}
	})
}

func newUploadRequest(t *testing.T, filename string, content string) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fileWriter.Write([]byte(content)); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
