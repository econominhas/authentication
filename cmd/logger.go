package main

import (
	"log/slog"
	"os"

	"github.com/econominhas/authentication/internal/models"
)

func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func newLogger() models.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     getLogLevel(),
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler)
}
