package handlers

import (
	"net/http"

	"logviewer/backend/parser"
	"logviewer/backend/store"
)

type uploadResponse struct {
	Count int            `json:"count"`
	Stats map[string]int `json:"stats"`
}

type UploadHandler struct {
	Parser parser.LogParser
	Store  store.Store
}

func (h UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusRequestEntityTooLarge, "file exceeds 10 MB")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing file field")
		return
	}
	defer file.Close()

	entries, err := h.Parser.ParseFile(file, header.Filename)
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, "parse failed: "+err.Error())
		return
	}

	for index := range entries {
		entries[index].Source = header.Filename
	}

	h.Store.Clear()
	h.Store.Add(entries)

	respondJSON(w, http.StatusOK, uploadResponse{
		Count: len(entries),
		Stats: h.Store.Stats(),
	})
}
