package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

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
	flag.Parse()

	// If no mdFilename as input, then show Usage:
	if *mdFilename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*mdFilename, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(mdFilename string, out io.Writer) error {
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
	fmt.Fprintln(out, outHtmlFileName)

	return saveHtml(outHtmlFileName, htmlContent)
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
