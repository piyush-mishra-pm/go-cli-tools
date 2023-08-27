package main

import (
	"flag"
	"fmt"
	"io"
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
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start crawling")
	// Action options
	listFiles := flag.Bool("list-file", false, "Only List files")
	// Filter options
	pickExtension := flag.String("pick-extension", "", "File extension to filter out")
	minSize := flag.Int64("min-size", 0, "Minimum file size")
	flag.Parse()

	launchConfigs := config{
		pickExtension: *pickExtension,
		minSize:       *minSize,
		listFiles:     *listFiles,
	}

	if err := run(*root, os.Stdout, launchConfigs); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(rootDir string, out io.Writer, launchConfig config) error {

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
			// List is the default option if nothing else was set
			return listFile(path, out)
		})
}
