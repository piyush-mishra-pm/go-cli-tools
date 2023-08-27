package main

import (
	"os"
	"testing"
)

func TestSkipFile(t *testing.T) {
	testCases := []struct {
		testName      string
		file          string
		pickExtension string
		minSize       int64
		expected      bool
	}{
		{testName: "FilterNoExtension", file: "testdata/dir.log", pickExtension: "", minSize: 0, expected: false},
		{testName: "FilterExtensionMatch", file: "testdata/dir.log", pickExtension: ".log", minSize: 0, expected: false},
		{testName: "FilterExtensionNoMatch", file: "testdata/dir.log", pickExtension: ".sh", minSize: 0, expected: true},
		{testName: "FilterExtensionSizeMatch", file: "testdata/dir.log", pickExtension: ".log", minSize: 10, expected: false},
		{testName: "FilterExtensionSizeNoMatch", file: "testdata/dir.log", pickExtension: ".log", minSize: 20, expected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			result := skipFile(tc.file, tc.pickExtension, tc.minSize, info)
			if result != tc.expected {
				t.Errorf("Expected '%t', got '%t' instead\n", tc.expected, result)
			}
		})
	}
}
