package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("Hi there, this string has seven words.")
	expectedWordCount := 7
	resWordCount := countWords(b)

	if expectedWordCount != resWordCount {
		t.Errorf("Expected word count %d, got %d instead.\n", expectedWordCount, resWordCount)
	}
}
