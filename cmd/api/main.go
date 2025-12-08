package main

import (
	"context"
	"net/http"
	"personalKnowledgeSearchEngine/internal/es"
	"personalKnowledgeSearchEngine/internal/notes"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	ctx := context.Background()
	esClient, err := es.NewESClient("https://localhost:9200")
	if err != nil {
		panic(err)
	}
	service := notes.NewService(ctx, esClient)
	handler := notes.NewHandler(service)

	r.Post("/notes", handler.CreateNote)
	r.Get("/notes/search", handler.SearchNotes)

	_ = http.ListenAndServe(":8080", r)
}
