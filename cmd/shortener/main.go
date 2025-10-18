package main

import (
	"github.com/Guram-Gurych/shortenerURL.git/internal/config"
	"github.com/Guram-Gurych/shortenerURL.git/internal/handler"
	"github.com/Guram-Gurych/shortenerURL.git/internal/logger"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	if err := logger.Initalize("info"); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	cfg := config.InitConfig()
	rep := repository.NewMemoryRepository()
	serv := service.NewShortenerService(rep)
	hndl := handler.NewHandler(serv, cfg.BaseURL)

	mux := chi.NewRouter()
	mux.Use(logger.RequestLogger)
	mux.Post("/", hndl.Post)
	mux.Get("/{id}", hndl.Get)
	mux.Post("/api/shorten", hndl.PostShorten)
	
	logger.Log.Info("Starting server", zap.String("address", cfg.ServerAddress))

	err := http.ListenAndServe(cfg.ServerAddress, mux)
	if err != nil {
		logger.Log.Fatal("Сервер упал", zap.Error(err))
	}
}
