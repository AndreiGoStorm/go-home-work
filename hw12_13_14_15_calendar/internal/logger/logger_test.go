package logger

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	const logFile = "log.txt"
	t.Run("log with info debug level", func(_ *testing.T) {
		msg := "log debug message"
		logger := New("Debug", logFile)

		logger.Info(msg)
		logger.Error(msg)
		logger.Warn(msg)
		logger.Debug(msg)

		logger.Close()

		data, err := os.ReadFile(logger.file.Name())
		require.NoError(t, err)

		err = os.Remove(logger.file.Name())
		require.NoError(t, err)

		arr := strings.Split(string(bytes.TrimSpace(data)), "\n")
		require.Equal(t, 4, len(arr))
	})

	t.Run("log with info level", func(_ *testing.T) {
		msg := "log info message"
		logger := New("Info", logFile)

		logger.Info(msg)
		logger.Error(msg)
		logger.Warn(msg)
		logger.Debug(msg)

		logger.Close()

		data, err := os.ReadFile(logger.file.Name())
		require.NoError(t, err)

		err = os.Remove(logger.file.Name())
		require.NoError(t, err)

		arr := strings.Split(string(bytes.TrimSpace(data)), "\n")
		require.Equal(t, 3, len(arr))
	})

	t.Run("log with warn level", func(_ *testing.T) {
		msg := "log warn message"
		logger := New("warn", logFile)

		logger.Info(msg)
		logger.Error(msg)
		logger.Warn(msg)
		logger.Debug(msg)

		logger.Close()

		data, err := os.ReadFile(logger.file.Name())
		require.NoError(t, err)

		err = os.Remove(logger.file.Name())
		require.NoError(t, err)

		arr := strings.Split(string(bytes.TrimSpace(data)), "\n")
		require.Equal(t, 2, len(arr))
	})

	t.Run("log with error level", func(_ *testing.T) {
		msg := "log error message"
		logger := New("error", logFile)

		logger.Info(msg)
		logger.Error(msg)
		logger.Warn(msg)
		logger.Debug(msg)

		logger.Close()

		data, err := os.ReadFile(logger.file.Name())
		require.NoError(t, err)

		err = os.Remove(logger.file.Name())
		require.NoError(t, err)

		arr := strings.Split(string(bytes.TrimSpace(data)), "\n")
		require.Equal(t, 1, len(arr))
	})

	t.Run("log change level", func(_ *testing.T) {
		msg := "log message"
		logger := New("info", logFile)

		logger.Info(msg)
		logger.Error(msg)
		logger.Warn(msg)
		logger.Debug(msg)

		data, err := os.ReadFile(logger.file.Name())
		require.NoError(t, err)

		arr := strings.Split(string(bytes.TrimSpace(data)), "\n")
		require.Equal(t, 3, len(arr))

		logger.SetLevel(slog.LevelWarn)

		logger.Info(msg)
		logger.Error(msg)
		logger.Warn(msg)
		logger.Debug(msg)

		logger.Close()

		data, err = os.ReadFile(logger.file.Name())
		require.NoError(t, err)

		err = os.Remove(logger.file.Name())
		require.NoError(t, err)

		arr = strings.Split(string(bytes.TrimSpace(data)), "\n")
		require.Equal(t, 5, len(arr))
	})
}
