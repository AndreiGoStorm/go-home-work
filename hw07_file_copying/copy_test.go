package main

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	tempPath      = "/tmp/"
	tempExtension = "*.txt"
)

func TestCopy(t *testing.T) {
	for _, test := range []struct {
		offset int64
		limit  int64
		file   string
	}{
		{offset: 0, limit: 0, file: "testdata/out_offset0_limit0.txt"},
		{offset: 0, limit: 10, file: "testdata/out_offset0_limit10.txt"},
		{offset: 0, limit: 1000, file: "testdata/out_offset0_limit1000.txt"},
		{offset: 0, limit: 10000, file: "testdata/out_offset0_limit10000.txt"},
		{offset: 100, limit: 1000, file: "testdata/out_offset100_limit1000.txt"},
		{offset: 6000, limit: 1000, file: "testdata/out_offset6000_limit1000.txt"},
	} {
		t.Run(fmt.Sprintf("Offset:%d Limit:%d", test.offset, test.limit), func(t *testing.T) {
			toFile, _ := os.CreateTemp(tempPath, tempExtension)
			defer os.Remove(toFile.Name())

			err := Copy("testdata/input.txt", toFile.Name(), test.offset, test.limit)
			require.NoError(t, err)

			fromFile, _ := os.Open(test.file)
			defer fromFile.Close()

			fromFileContent, _ := os.ReadFile(fromFile.Name())
			toFileContent, _ := os.ReadFile(toFile.Name())
			require.Equal(t, fromFileContent, toFileContent)
		})
	}

	t.Run("Error: not existed file", func(t *testing.T) {
		err := Copy("input.txt", "output.txt", 0, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("Error: file is directory", func(t *testing.T) {
		err := Copy(tempPath, tempPath+"output.txt", 0, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("Error: device file", func(t *testing.T) {
		toFile, _ := os.CreateTemp(tempPath, tempExtension)
		defer os.Remove(toFile.Name())

		err := Copy("/dev/urandom", toFile.Name(), 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("Error: offset more then size", func(t *testing.T) {
		toFile, _ := os.CreateTemp(tempPath, tempExtension)
		defer os.Remove(toFile.Name())

		err := Copy(toFile.Name(), tempPath+"output.txt", 1000, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})
}
