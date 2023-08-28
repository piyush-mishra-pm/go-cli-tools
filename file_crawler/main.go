package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	// file extensions to consider.
	pickExtension string
	// min file size to consider.
	minSize int64
	// list files?
	listFiles bool
	// Delete files"
	deleteFiles bool
	// log File
	logFile io.Writer
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start crawling")
	logFile := flag.String("log", "", "Logs delete ops")
	// Action options
	listFiles := flag.Bool("list-file", false, "Only List files")
	deleteFiles := flag.Bool("delete-file", false, "Delete files")
	// Filter options
	pickExtension := flag.String("pick-extension", "", "File extension to filter out")
	minSize := flag.Int64("min-size", 0, "Minimum file size")
	flag.Parse()

	// Open log file (if log file provided)
	var (
		file = os.Stdout
		err  error
	)
	if *logFile != "" {
		file, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
	}

	launchConfigs := config{
		pickExtension: *pickExtension,
		minSize:       *minSize,
		listFiles:     *listFiles,
		deleteFiles:   *deleteFiles,
		logFile:       file,
	}

	if err := run(*root, os.Stdout, launchConfigs); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(rootDir string, out io.Writer, launchConfig config) error {

	onDeleteLoggerFunc := log.New(launchConfig.logFile, "DELETED_FILE:", log.LstdFlags)

	return filepath.Walk(rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if skipFile(path, launchConfig.pickExtension, launchConfig.minSize, info) {
				return nil
			}
			// If list was explicitly set, don't do anything else
			if launchConfig.listFiles {
				return listFile(path, out)
			}

			// Delete Files:
			if launchConfig.deleteFiles {
				return deleteFile(path, onDeleteLoggerFunc)
			}

			// List is the default option if nothing else was set
			return listFile(path, out)
		})
}
