package logger

import (
	"email-verification-service/internal/pkg/config"
	"email-verification-service/utils"
	"io"
	"log"
	"log/slog"
	"os"
)

type Logger = *slog.Logger

func New(cfg config.Logger) *Logger {
	logLevel := parseLogLevel(cfg.LogLevel)

	options := &slog.HandlerOptions{
		Level: logLevel,
	}

	if logLevel == slog.LevelDebug {
		options.AddSource = true
	}

	file, err := utils.FileHandle(cfg.FilePath)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}

	w := io.MultiWriter(file, os.Stdout)
	logger := slog.New(slog.NewJSONHandler(w, options))
	slog.SetDefault(logger)

	return &logger
}

func parseLogLevel(level string) slog.Level {
	logLevelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	logLevel, ok := logLevelMap[level]

	if !ok {
		logLevel = slog.LevelDebug
	}

	return logLevel
}
