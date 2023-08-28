package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
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

func deleteFile(path string, loggerFunc *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	loggerFunc.Println(path)
	return nil
}

func archiveFile(destDir, root, fileOriginalPath string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}
	relDir, err := filepath.Rel(root, filepath.Dir(fileOriginalPath))
	if err != nil {
		return err
	}
	destFileBaseName := fmt.Sprintf("%s.gz", filepath.Base(fileOriginalPath))
	targetPath := filepath.Join(destDir, relDir, destFileBaseName)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	newArchivedFile, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer newArchivedFile.Close()

	originalFile, err := os.Open(fileOriginalPath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	gzipWriter := gzip.NewWriter(newArchivedFile)
	gzipWriter.Name = filepath.Base(fileOriginalPath)
	// not deferring the call to zw.Close() to ensure we return
	// any potential errors because, if the compressing fails,
	// the calling function will get an error and decide how to proceed.

	if _, err = io.Copy(gzipWriter, originalFile); err != nil {
		return err
	}
	if err := gzipWriter.Close(); err != nil {
		return err
	}

	return newArchivedFile.Close()
}
