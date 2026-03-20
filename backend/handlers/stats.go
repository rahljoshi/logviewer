package handlers

import (
	"net/http"

	"logviewer/backend/store"
)

type statsResponse struct {
	Total  int            `json:"total"`
	Levels map[string]int `json:"levels"`
}

type StatsHandler struct {
	Store store.Store
}

func (h StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	levels := h.Store.Stats()
	total := 0
	for _, count := range levels {
		total += count
	}

	respondJSON(w, http.StatusOK, statsResponse{
		Total:  total,
		Levels: levels,
	})
}
