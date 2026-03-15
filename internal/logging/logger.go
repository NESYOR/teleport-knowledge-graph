package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// New builds a structured logger configured by level and format.
func New(level string, format string) *slog.Logger {
	lvl := parseLevel(level)
	opts := &slog.HandlerOptions{Level: lvl}

	var handler slog.Handler
	if strings.EqualFold(format, "text") {
		handler = slog.NewTextHandler(output(), opts)
	} else {
		handler = slog.NewJSONHandler(output(), opts)
	}
	return slog.New(handler)
}

func output() io.Writer {
	return os.Stdout
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
