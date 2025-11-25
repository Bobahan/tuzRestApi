package logger

import (
	"log/slog"
	"os"
)

const (
	LOCAL       = "local"
	DEVELOPMENT = "development"
	PRODUCTION  = "production"
)

func New(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case LOCAL:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case DEVELOPMENT:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case PRODUCTION:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
