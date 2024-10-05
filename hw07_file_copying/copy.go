package main

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	bytes, err := validateFile(fromPath, offset, limit)
	if err != nil {
		return err
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	if offset > 0 {
		fromFile.Seek(offset, 0)
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	bar := pb.Full.Start64(bytes)
	toFileWriter := bar.NewProxyWriter(toFile)
	defer bar.Finish()

	_, err = io.CopyN(toFileWriter, fromFile, bytes)
	if err != nil {
		return err
	}

	return nil
}

func validateFile(name string, offset, limit int64) (int64, error) {
	info, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fs.ErrNotExist
		}
		return 0, err
	}
	if isUnsupportedFile(info.Mode()) {
		return 0, ErrUnsupportedFile
	}
	if offset > info.Size() {
		return 0, ErrOffsetExceedsFileSize
	}
	bytes := info.Size() - offset
	if limit > 0 && bytes > limit {
		bytes = limit
	}
	return bytes, nil
}

func isUnsupportedFile(m fs.FileMode) bool {
	return m&os.ModeDevice != 0 || m&os.ModeCharDevice != 0 || m&os.ModeDir != 0 || m&os.ModeSocket != 0
}
