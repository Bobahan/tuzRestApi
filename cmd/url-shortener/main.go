package main

import (
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// config +
	// logger +
	// storage +
	// router +
	// server +

	config := config.New()

	logger := logger.New(config.Env)

	logger.Info("server has been started")
	logger.Debug("debug logger")

	storage, err := sqlite.New(config.StoragePath)

	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	if err != nil {
		logger.Error("failed to get url", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			config.HttpServer.User: config.HttpServer.Password,
		}))

		r.Post("/", url.New(logger, storage))
	})

	router.Get("/{alias}", redirect.New(logger, storage))

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HttpServer.Timeout,
		WriteTimeout: config.HttpServer.Timeout,
		IdleTimeout:  config.HttpServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}
}
