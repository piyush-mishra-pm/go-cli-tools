package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	projDir := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*projDir, os.Stdout); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func run(projDir string, out io.Writer) error {
	if projDir == "" {
		return fmt.Errorf("❗️ Missing Project directory: %w", ErrValidation)
	}

	args := []string{"build", ".", "errors"}

	cmd := exec.Command("go", args...)

	cmd.Dir = projDir

	if err := cmd.Run(); err != nil {
		return &errStep{step: "go build", msg: "❗️ 'go build' failed", cause: err}
	}

	_, err := fmt.Fprintln(out, "✅ Go Build: SUCCESS")

	return err
}
