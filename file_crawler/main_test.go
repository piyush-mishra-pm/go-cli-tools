package main

import (
	"bytes"
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
