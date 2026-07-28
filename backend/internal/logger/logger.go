package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

func Init(level string) {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}),
	)

	slog.SetDefault(logger)
}

func DebugEnabled() bool {
	return slog.Default().Enabled(context.Background(), slog.LevelDebug)
}
