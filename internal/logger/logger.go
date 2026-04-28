package logger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
