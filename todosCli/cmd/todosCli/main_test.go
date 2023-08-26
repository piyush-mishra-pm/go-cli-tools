package main_test

import (
	"fmt"
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
	newTodoTaskName := "test task number 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	binFilePath := filepath.Join(dir, tempBinFileName)

	t.Run("AddNewTodo", func(t *testing.T) {
		cmd := exec.Command(binFilePath, "-add", newTodoTaskName)
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

		expected := fmt.Sprintf("[ ] 1: %s\n", newTodoTaskName)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
