package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	logLevel *slog.LevelVar
	slog     *slog.Logger
	file     *os.File
}

func New(level, logFile string) *Logger {
	logger := new(Logger)
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return nil
	}
	logger.logLevel = &slog.LevelVar{}
	logger.slog = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
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
	logger.file = file

	return logger
}

func (l *Logger) SetLevel(level slog.Level) {
	l.logLevel.Set(level)
}

func (l *Logger) Info(msg string) {
	l.slog.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.slog.Error(msg)
}

func (l *Logger) Warn(msg string) {
	l.slog.Warn(msg)
}

func (l *Logger) Debug(msg string) {
	l.slog.Debug(msg)
}

func (l *Logger) Close() {
	l.file.Close()
}
