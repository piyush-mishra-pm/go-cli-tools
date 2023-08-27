package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	goldenInputMdFile = "./testdata/test1.md"
	outputHtmlFile    = "test1.md.html"
	goldenHtmlFile    = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	goldenMdFileContent, err := os.ReadFile(goldenInputMdFile)
	if err != nil {
		t.Fatal(err)
	}
	resultHtmlContent := parseMdFileContent(goldenMdFileContent)
	expectedHtmlContent, err := os.ReadFile(goldenHtmlFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expectedHtmlContent, resultHtmlContent) {
		t.Logf("goldenHtml:\n%s\n", expectedHtmlContent)
		t.Logf("resultHtml:\n%s\n", resultHtmlContent)
		t.Error("Result html content does not match golden html file")
	}
}

func TestRun(t *testing.T) {
	if err := run(goldenInputMdFile); err != nil {
		t.Fatal(err)
	}
	resultHtmlContent, err := os.ReadFile(outputHtmlFile)
	if err != nil {
		t.Fatal(err)
	}
	expectedHtmlContent, err := os.ReadFile(goldenHtmlFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expectedHtmlContent, resultHtmlContent) {
		t.Logf("goldenHtml:\n%s\n", expectedHtmlContent)
		t.Logf("resultHtml:\n%s\n", resultHtmlContent)
		t.Error("Result html content does not match golden html file")
	}
	os.Remove(outputHtmlFile)
}
