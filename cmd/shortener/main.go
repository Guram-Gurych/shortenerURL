package main

import (
	"github.com/Guram-Gurych/shortenerURL.git/internal/config"
	"github.com/Guram-Gurych/shortenerURL.git/internal/config/db"
	"github.com/Guram-Gurych/shortenerURL.git/internal/handler"
	"github.com/Guram-Gurych/shortenerURL.git/internal/logger"
	"github.com/Guram-Gurych/shortenerURL.git/internal/middleware"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	cfg := config.InitConfig()
	if cfg.DatabaseDSN != "" {
		db, err := db.Initialize(cfg.DatabaseDSN)
		if err != nil {
			logger.Log.Fatal("Ошибка инициализации DB", zap.Error(err))
		}
		defer db.Close()
	}

	rep, err := repository.NewFileRepository(cfg.FileStoragePath)
	if err != nil {
		logger.Log.Fatal("Ошибка репозитория", zap.Error(err))
	}
	defer rep.Close()

	serv := service.NewShortenerService(rep)
	hndl := handler.NewHandler(serv, cfg.BaseURL, db)

	mux := chi.NewRouter()
	mux.Use(middleware.RequestLogger)
	mux.Use(middleware.GzipMiddleware)
	mux.Post("/", hndl.Post)
	mux.Get("/{id}", hndl.Get)
	mux.Post("/api/shorten", hndl.PostShorten)
	mux.Get("/ping", hndl.GetPing)

	logger.Log.Info("Starting server", zap.String("address", cfg.ServerAddress))

	err = http.ListenAndServe(cfg.ServerAddress, mux)
	if err != nil {
		logger.Log.Fatal("Сервер упал", zap.Error(err))
	}
}
