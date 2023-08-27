package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	htmlHeader = `<!DOCTYPE html>
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>Markdown Preview Tool</title>
	</head>
	<body>
	`
	htmlFooter = `
	</body>
	</html>
	`
)

func main() {
	// Parse flags:
	mdFilename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("skip-preview", false, "Skip automatic Markdown Preview")
	flag.Parse()

	// If no mdFilename as input, then show Usage:
	if *mdFilename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*mdFilename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(mdFilename string, out io.Writer, skipPreview bool) error {
	mdFileContent, err := os.ReadFile(mdFilename)
	if err != nil {
		return err
	}

	htmlContent := parseMdFileContent(mdFileContent)

	// Write in temp file and check for errors.
	tempFile, err := os.CreateTemp("", "markdown_preview*.html")
	if err != nil {
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}

	outHtmlFileName := tempFile.Name()
	defer os.Remove(outHtmlFileName)

	fmt.Fprintln(out, outHtmlFileName)

	if err := saveHtml(outHtmlFileName, htmlContent); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}
	return preview(outHtmlFileName)
}

func parseMdFileContent(mdFileContent []byte) []byte {
	// Parse and Sanitise HTML Body using md file content.
	htmlBodyOutput := blackfriday.Run(mdFileContent)
	sanitisedHtmlBodyOutput := bluemonday.UGCPolicy().SanitizeBytes(htmlBodyOutput)

	// Add header and Footer to Body.
	var buffer bytes.Buffer
	buffer.WriteString(htmlHeader)
	buffer.Write(sanitisedHtmlBodyOutput)
	buffer.WriteString(htmlFooter)

	return buffer.Bytes()
}

func saveHtml(outHtmlFileName string, htmlContent []byte) error {
	return os.WriteFile(outHtmlFileName, htmlContent, 0644)
}

func preview(fname string) error {
	cName := ""
	cmdParams := []string{}
	// Define execute utility based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cmdParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	// Append filename to command parameters
	cmdParams = append(cmdParams, fname)
	// Locate executable in PATH
	cmdPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using default program
	err = exec.Command(cmdPath, cmdParams...).Run()

	// Wait for browser to open file before (automatically) deleting.
	time.Sleep(3 * time.Second)
	return err
}
