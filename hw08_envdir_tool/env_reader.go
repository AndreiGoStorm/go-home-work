package main

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(files))
	for _, file := range files {
		if isUnsupportedFile(file.Type()) {
			continue
		}
		if strings.Contains(file.Name(), "=") {
			continue
		}
		line, err := readFirstLineFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		env[file.Name()] = *formatLine(line)
	}
	return env, nil
}

func readFirstLineFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	line, _, err := bufio.NewReader(file).ReadLine()
	if err == io.EOF { //nolint:errorlint
		return []byte(""), nil
	}
	if err != nil {
		return nil, err
	}
	return line, nil
}

func formatLine(line []byte) *EnvValue {
	formated := bytes.TrimRight(line, " \t")
	formated = bytes.ReplaceAll(formated, []byte("\x00"), []byte("\n"))
	if len(formated) == 0 {
		return &EnvValue{NeedRemove: true}
	}
	return &EnvValue{Value: string(formated)}
}

func isUnsupportedFile(m fs.FileMode) bool {
	return m&os.ModeDevice != 0 || m&os.ModeCharDevice != 0 || m&os.ModeDir != 0 || m&os.ModeSocket != 0
}
