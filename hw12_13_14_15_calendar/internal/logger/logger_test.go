package logger

import (
	"errors"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("log with info debug level", func(_ *testing.T) {
		infoMsg := "log debug message"
		err := errors.New("error message")
		logger := New("Debug")

		logger.Info(infoMsg)
		logger.Error("error in log", err)
		logger.Warn("warn in log", err)
		logger.Debug("debug in log", err)

		/*
			time=2024-11-12T17:47:55.741+01:00 level=INFO msg="log debug message"
			time=2024-11-12T17:47:55.741+01:00 level=ERROR msg="error in log" error="error message"
			time=2024-11-12T17:47:55.741+01:00 level=WARN msg="warn in log" warning="error message"
			time=2024-11-12T17:47:55.741+01:00 level=DEBUG msg="debug in log" debug="error message"
		*/
	})

	t.Run("log with info level", func(_ *testing.T) {
		infoMsg := "log info message"
		err := errors.New("error message")
		logger := New("Info")

		logger.Info(infoMsg)
		logger.Error("error in log", err)
		logger.Warn("warn in log", err)
		logger.Debug("debug in log", err)

		/*
			time=2024-11-12T17:48:35.902+01:00 level=INFO msg="log debug message"
			time=2024-11-12T17:48:35.902+01:00 level=ERROR msg="error in log" error="error message"
			time=2024-11-12T17:48:35.902+01:00 level=WARN msg="warn in log" warning="error message"
		*/
	})

	t.Run("log with warn level", func(_ *testing.T) {
		infoMsg := "log warn message"
		err := errors.New("error message")
		logger := New("warn")

		logger.Info(infoMsg)
		logger.Error("error in log", err)
		logger.Warn("warn in log", err)
		logger.Debug("debug in log", err)

		/*
			time=2024-11-12T17:49:30.008+01:00 level=ERROR msg="error in log" error="error message"
			time=2024-11-12T17:49:30.008+01:00 level=WARN msg="warn in log" warning="error message"
		*/
	})

	t.Run("log with error level", func(_ *testing.T) {
		infoMsg := "log error message"
		err := errors.New("error message")
		logger := New("error")

		logger.Info(infoMsg)
		logger.Error("error in log", err)
		logger.Warn("warn in log", err)
		logger.Debug("debug in log", err)

		/*
			time=2024-11-12T17:49:45.802+01:00 level=ERROR msg="error in log" error="error message"
		*/
	})
}
