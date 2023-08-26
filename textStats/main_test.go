package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("Hi there, this string has seven words.")
	expectedWordCount := 7
	resWordCount := countWords(b, false)

	if expectedWordCount != resWordCount {
		t.Errorf("Expected word count %d, got %d instead.\n", expectedWordCount, resWordCount)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("Hi\nHave 3\nlines.")
	expLineCount := 3
	resLineCount := countWords(b, true)

	if expLineCount != resLineCount {
		t.Errorf("Expected line count %d, got %d instead.\n", expLineCount, resLineCount)
	}
}
