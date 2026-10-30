package logger

import (
	"log/slog"
	"os"

	"github.com/top-coins/internal/models"
)

type Logger struct {
	*slog.Logger
}

func NewSlogLogger(env models.Environment) *Logger {
	var logHanlder slog.Handler

	switch env {
	case models.Dev:
		logHanlder = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case models.Debug:
		logHanlder = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	default:
		logHanlder = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return &Logger{slog.New(logHanlder)}
}
