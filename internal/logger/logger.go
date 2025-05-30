package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
)

func New(env string) *slog.Logger {

	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return logger
}
