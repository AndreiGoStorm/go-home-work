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
	ErrSameFilePaths         = errors.New("from and to file paths the same")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	bytes, err := validate(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	if offset > 0 {
		fromFile.Seek(offset, 0)
	}

	bar := pb.Full.Start64(bytes)
	toFileWriter := bar.NewProxyWriter(toFile)
	defer bar.Finish()

	_, err = io.CopyN(toFileWriter, fromFile, bytes)
	if err != nil {
		return err
	}

	return nil
}

func validate(from, to string, offset, limit int64) (int64, error) {
	infoFrom, _ := os.Stat(from)
	if isUnsupportedFile(infoFrom.Mode()) {
		return 0, ErrUnsupportedFile
	}
	if offset > infoFrom.Size() {
		return 0, ErrOffsetExceedsFileSize
	}
	infoTo, _ := os.Stat(to)
	if os.SameFile(infoFrom, infoTo) {
		return 0, ErrSameFilePaths
	}
	bytes := infoFrom.Size() - offset
	if limit > 0 && bytes > limit {
		bytes = limit
	}
	return bytes, nil
}

func isUnsupportedFile(m fs.FileMode) bool {
	return m&os.ModeDevice != 0 || m&os.ModeCharDevice != 0 || m&os.ModeDir != 0 || m&os.ModeSocket != 0
}
