package main

import (
	"github.com/Guram-Gurych/shortenerURL.git/internal/config"
	"github.com/Guram-Gurych/shortenerURL.git/internal/handler"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg := config.InitConfig()
	rep := repository.NewMemoryRepository()
	serv := service.NewShortenerService(rep)
	hndl := handler.NewHandler(serv, cfg.BaseURL)

	mux := chi.NewRouter()
	mux.Post("/", hndl.Post)
	mux.Get("/{id}", hndl.Get)

	err := http.ListenAndServe(cfg.ServerAddress, mux)
	if err != nil {
		log.Fatalf("Сервер упал: %v", err)
	}
}
