package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	tempBinFileName = "todo"
	saveFilename    = ".todo.json"
)

func TestMain(m *testing.M) {
	// Temporarily disable filename set in envt variable.
	val := os.Getenv("TODO_FILENAME")
	if val != "" {
		os.Unsetenv("TODO_FILENAME")
		defer os.Setenv("TODO_FILENAME", val)
	}

	fmt.Println("🔧 Building todosCli tool :")

	if runtime.GOOS == "windows" {
		tempBinFileName += ".exe"
	}

	build := exec.Command("go", "build", "-o", tempBinFileName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error in building todosCli tool: %s", err)
		os.Exit(1)
	}

	fmt.Println("✅ Build successful.\n🔍Running tests ...")
	testResults := m.Run()

	fmt.Println("✅ Ran tests.\n🧹Cleaning up.")
	os.Remove(tempBinFileName)
	os.Remove(saveFilename)

	os.Exit(testResults)
}

func TestTodoCLI(t *testing.T) {
	newTodoTaskName1 := "test task number 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	binFilePath := filepath.Join(dir, tempBinFileName)

	t.Run("AddNewTodoFromArguments", func(t *testing.T) {
		cmd := exec.Command(binFilePath, "-add", newTodoTaskName1)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	newTodoTaskName2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(binFilePath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(cmdStdIn, newTodoTaskName2); err != nil {
			t.Fatal(err)
		}
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTodos", func(t *testing.T) {
		cmd := exec.Command(binFilePath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("[ ] 1: %s\n[ ] 2: %s\n", newTodoTaskName1, newTodoTaskName2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
