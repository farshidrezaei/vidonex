// Package main implements the Vidonyx command-line video composition interface.
package main

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// SetupLogger initializes and returns a structured slog.Logger based on log level and format.
func SetupLogger(levelName, format string, destination io.Writer) *slog.Logger {
	if destination == nil {
		destination = os.Stderr
	}

	var level slog.Level
	switch strings.ToLower(strings.TrimSpace(levelName)) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "info":
		fallthrough
	default:
		level = slog.LevelInfo
	}

	options := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(strings.TrimSpace(format)) == "json" {
		handler = slog.NewJSONHandler(destination, options)
	} else {
		handler = slog.NewTextHandler(destination, options)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
