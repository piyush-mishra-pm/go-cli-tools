package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name          string
		root          string
		launchConfigs config
		expected      string
	}{
		{name: "NoFilter", root: "testdata",
			launchConfigs: config{pickExtension: "", minSize: 0, listFiles: true},
			expected:      "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata",
			launchConfigs: config{pickExtension: ".log", minSize: 0, listFiles: true},
			expected:      "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata",
			launchConfigs: config{pickExtension: ".log", minSize: 10, listFiles: true},
			expected:      "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata",
			launchConfigs: config{pickExtension: ".log", minSize: 20, listFiles: true},
			expected:      ""},
		{name: "FilterExtensionNoMatch", root: "testdata",
			launchConfigs: config{pickExtension: ".gz", minSize: 0, listFiles: true},
			expected:      ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := run(tc.root, &buffer, tc.launchConfigs); err != nil {
				t.Fatal(err)
			}
			res := buffer.String()
			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name          string
		launchConfigs config
		extNoDelete   string
		nDelete       int
		nNoDelete     int
		expected      string
	}{
		{name: "DeleteExtensionNoMatch",
			launchConfigs: config{pickExtension: ".log", deleteFiles: true},
			extNoDelete:   ".gz", nDelete: 0, nNoDelete: 10,
			expected: ""},
		{name: "DeleteExtensionMatch",
			launchConfigs: config{pickExtension: ".log", deleteFiles: true},
			extNoDelete:   "", nDelete: 10, nNoDelete: 0,
			expected: ""},
		{name: "DeleteExtensionMixed",
			launchConfigs: config{pickExtension: ".log", deleteFiles: true},
			extNoDelete:   ".gz", nDelete: 5, nNoDelete: 5,
			expected: ""},
	}
	// Execute RunDel test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)
			tc.launchConfigs.logFile = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.launchConfigs.pickExtension: tc.nDelete,
				tc.extNoDelete:                 tc.nNoDelete,
			})
			defer cleanup()
			if err := run(tempDir, &buffer, tc.launchConfigs); err != nil {
				t.Fatal(err)
			}
			res := buffer.String()
			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}
			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}
			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d files left, got %d instead\n",
					tc.nNoDelete, len(filesLeft))
			}
			// Verify log output (expected number of log lines must have been created)
			expLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n",
					expLogLines, len(lines))
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "file_crawler_test")
	if err != nil {
		t.Fatal(err)
	}
	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}
