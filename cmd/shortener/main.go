package main

import (
	"github.com/Guram-Gurych/shortenerURL.git/internal/handler"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	var baseURL = "http://localhost:8080"
	var serverAddr = ":8080"

	repository := repository.NewMemoryRepository()
	server := service.NewShortenerService(repository)
	handle := handler.NewHandler(server, baseURL)

	mux := chi.NewRouter()
	mux.Post("/", handle.Post)
	mux.Get("/{id}", handle.Get)

	err := http.ListenAndServe(serverAddr, mux)
	if err != nil {
		log.Fatalf("Сервер упал: %v", err)
	}
}
