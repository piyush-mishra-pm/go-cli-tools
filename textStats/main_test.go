package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("Hi there, this string has seven words.")
	expectedWordCount := 7
	resWordCount := count(b, false, false)

	if expectedWordCount != resWordCount {
		t.Errorf("Expected word count %d, got %d instead.\n", expectedWordCount, resWordCount)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("Hi\nHave 3\nlines.")
	expLineCount := 3
	resLineCount := count(b, true, false)

	if expLineCount != resLineCount {
		t.Errorf("Expected line count %d, got %d instead.\n", expLineCount, resLineCount)
	}
}

func TestCountBytes(t *testing.T) {
	b := bytes.NewBufferString("SWIFT")
	// byte array: [83 87 73 70 84]
	expByteCount := 5
	resByteCount := count(b, false, true)

	if expByteCount != resByteCount {
		t.Errorf("Expected byte count %d, got %d instead.\n", expByteCount, resByteCount)
	}
}
