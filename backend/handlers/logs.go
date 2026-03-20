package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"logviewer/backend/models"
	"logviewer/backend/store"
)

type logsResponse struct {
	Entries  []models.LogEntry `json:"entries"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

type okResponse struct {
	OK bool `json:"ok"`
}

type LogsHandler struct {
	Store store.Store
}

func (h LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodDelete:
		h.Store.Clear()
		respondJSON(w, http.StatusOK, okResponse{OK: true})
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h LogsHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	page, err := parsePositiveInt(r.URL.Query().Get("page"), 1)
	if err != nil {
		respondError(w, http.StatusBadRequest, "page must be greater than or equal to 1")
		return
	}

	pageSize, err := parsePositiveInt(r.URL.Query().Get("pageSize"), 100)
	if err != nil || pageSize > 500 {
		respondError(w, http.StatusBadRequest, "pageSize must be between 1 and 500")
		return
	}

	opts := store.FilterOpts{
		Query:    r.URL.Query().Get("q"),
		Level:    r.URL.Query().Get("level"),
		Page:     page,
		PageSize: pageSize,
	}

	entries, total := h.Store.Filter(opts)
	respondJSON(w, http.StatusOK, logsResponse{
		Entries:  entries,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func parsePositiveInt(raw string, defaultValue int) (int, error) {
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	if value < 1 {
		return 0, fmt.Errorf("value must be >= 1")
	}

	return value, nil
}
