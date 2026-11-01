package logger

import (
	"log/slog"
	"os"

	"github.com/top-coins/pkg/models"
)

type CustomLogger struct {
	*slog.Logger
}

func NewSlogLogger(env models.Environment) *CustomLogger {
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

	return &CustomLogger{slog.New(logHanlder)}
}

func (l *CustomLogger) Attr(key string, value any) slog.Attr {
	val := slog.AnyValue(value)

	if err, ok := value.(error); ok {
		val = slog.AnyValue(err.Error())
	}

	return slog.Attr{
		Key:   key,
		Value: val,
	}
}
