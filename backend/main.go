package main

import (
	"log"
	"net/http"

	"logviewer/backend/handlers"
	"logviewer/backend/middleware"
	"logviewer/backend/parser"
	"logviewer/backend/store"
)

func main() {
	log.SetFlags(0)

	memStore := store.NewMemStore()
	logParser := parser.NewAutoParser()

	mux := http.NewServeMux()
	mux.Handle("/api/upload", handlers.UploadHandler{
		Parser: logParser,
		Store:  memStore,
	})
	mux.Handle("/api/logs", handlers.LogsHandler{
		Store: memStore,
	})
	mux.Handle("/api/stats", handlers.StatsHandler{
		Store: memStore,
	})
	mux.Handle("/api/health", handlers.HealthHandler{})

	handler := middleware.CORS(middleware.Logger(mux))

	log.Println("LogViewer backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
