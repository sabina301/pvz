package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

func Init(level string) {
	var slogLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	filePath := os.Getenv("LOG_FILE")
	if filePath == "" {
		filePath = "default.log"
	}

	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})
	fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slogLevel,
	})

	multiHandler := NewMultiHandler(stdoutHandler, fileHandler)

	Log = slog.New(multiHandler)
}
