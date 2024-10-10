package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Correct reading temp dir", func(t *testing.T) {
		testDir, _ := os.MkdirTemp("", "testdir")
		defer os.Remove(testDir)
		toFile, _ := os.CreateTemp(testDir, "")
		defer os.Remove(toFile.Name())
		value := "first line"
		_, _ = toFile.WriteString(value)

		expected := Environment{filepath.Base(toFile.Name()): EnvValue{Value: value, NeedRemove: false}}

		env, err := ReadDir(testDir)

		require.Equal(t, len(env), 1)
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("Correct reading testdata dir", func(t *testing.T) {
		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		env, err := ReadDir("testdata/env")

		require.Equal(t, len(env), 5)
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("Error: directory is not existed", func(t *testing.T) {
		env, err := ReadDir("testdata/not_existed/")

		require.Len(t, env, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("Error: directory is empty", func(t *testing.T) {
		testDir, _ := os.MkdirTemp("", "empty")
		defer os.Remove(testDir)

		env, err := ReadDir(testDir)

		require.Len(t, env, 0)
		require.NoError(t, err)
	})

	t.Run("Error: unsupported files", func(t *testing.T) {
		testDir, _ := os.MkdirTemp("", "unsupported")
		defer os.Remove(testDir)
		dirTemp, _ := os.MkdirTemp(testDir, "temp")
		defer os.Remove(dirTemp)
		toFile, _ := os.CreateTemp(testDir, "22=1.txt")
		defer os.Remove(toFile.Name())

		env, err := ReadDir(testDir)

		require.Len(t, env, 0)
		require.NoError(t, err)
	})
}
