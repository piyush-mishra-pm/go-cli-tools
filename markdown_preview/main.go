package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
	{{ .Body }}
	</body>
</html>
`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	// Parse flags:
	mdFilename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("skip-preview", false, "Skip automatic Markdown Preview")
	templateFileName := flag.String("template", "", "User defined Template filename")
	flag.Parse()

	// If no mdFilename as input, then show Usage:
	if *mdFilename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*mdFilename, *templateFileName, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(mdFilename, tempFileName string, out io.Writer, skipPreview bool) error {
	mdFileContent, err := os.ReadFile(mdFilename)
	if err != nil {
		return err
	}

	htmlContent, err := parseMdFileContent(mdFileContent, tempFileName)
	if err != nil {
		return err
	}

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

	if err := saveHtml(outHtmlFileName, htmlContent); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outHtmlFileName)

	return preview(outHtmlFileName)
}

func parseMdFileContent(mdFileContent []byte, templateFileName string) ([]byte, error) {
	// Parse and Sanitise HTML Body using md file content.
	htmlBodyOutput := blackfriday.Run(mdFileContent)
	sanitisedHtmlBodyOutput := bluemonday.UGCPolicy().SanitizeBytes(htmlBodyOutput)

	// Parse the contents of the defaultTemplate const into a new Template
	t, err := template.New("markdown_preview").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}
	// If user provided alternate template file, replace template
	if templateFileName != "" {
		t, err = template.ParseFiles(templateFileName)
		if err != nil {
			return nil, err
		}
	}
	// instantiate the content struct, adding the title and body
	contentData := content{
		Title: "Markdown Preview Go Tool",
		Body:  template.HTML(sanitisedHtmlBodyOutput),
	}

	// Add header and Footer to Body.
	var buffer bytes.Buffer
	// Execute the template with the content type
	if err := t.Execute(&buffer, contentData); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
