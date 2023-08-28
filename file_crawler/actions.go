package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func skipFile(path, pickExtension string, minSize int64, fileInfo os.FileInfo) bool {
	if fileInfo.IsDir() || fileInfo.Size() < minSize {
		return true
	}
	if pickExtension != "" && filepath.Ext(path) != pickExtension {
		return true
	}
	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func deleteFile(path string) error {
	return os.Remove(path)
}
