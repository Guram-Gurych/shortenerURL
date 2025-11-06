package main

import (
	"database/sql"
	"github.com/Guram-Gurych/shortenerURL.git/internal/config"
	"github.com/Guram-Gurych/shortenerURL.git/internal/config/db"
	"github.com/Guram-Gurych/shortenerURL.git/internal/handler"
	"github.com/Guram-Gurych/shortenerURL.git/internal/logger"
	"github.com/Guram-Gurych/shortenerURL.git/internal/middleware"
	"github.com/Guram-Gurych/shortenerURL.git/internal/repository"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	cfg := config.InitConfig()

	var dbConn *sql.DB
	var err error
	if cfg.DatabaseDSN != "" {
		dbConn, err = db.Initialize(cfg.DatabaseDSN)
		if err != nil {
			logger.Log.Fatal("Ошибка инициализации DB", zap.Error(err))
		}
		defer dbConn.Close()
		logger.Log.Info("DB connection established")
	}

	rep, err := repository.NewFileRepository(cfg.FileStoragePath)
	if err != nil {
		logger.Log.Fatal("Ошибка репозитория", zap.Error(err))
	}
	defer rep.Close()

	serv := service.NewShortenerService(rep)
	hndl := handler.NewHandler(serv, cfg.BaseURL, dbConn)

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
