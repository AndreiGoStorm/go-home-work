package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	logLevel *slog.LevelVar
	Slog     *slog.Logger
}

func New(level string) *Logger {
	logger := new(Logger)
	logger.logLevel = &slog.LevelVar{}
	logger.Slog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logger.logLevel,
	}))

	var l slog.Level
	switch strings.ToLower(level) {
	case "error":
		l = slog.LevelError
	case "warn":
		l = slog.LevelWarn
	case "debug":
		l = slog.LevelDebug
	default:
		l = slog.LevelInfo
	}
	logger.SetLevel(l)

	return logger
}

func (l *Logger) SetLevel(level slog.Level) {
	l.logLevel.Set(level)
}

func (l *Logger) Info(msg string) {
	l.Slog.Info(msg)
}

func (l *Logger) Error(msg string, err error) {
	l.Slog.Error(msg, "error", err)
}

func (l *Logger) Warn(msg string, err error) {
	l.Slog.Warn(msg, "warning", err)
}

func (l *Logger) Debug(msg string, err error) {
	l.Slog.Debug(msg, "debug", err)
}
