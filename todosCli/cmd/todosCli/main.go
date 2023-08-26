package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	todoscli "piyush-mishra-pm/go-cli-tools/todosCli"
	"strings"
)

// Hardcoded save file name:
var todoListFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed by Piyush Mishra\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parse Flags:
	addTodo := flag.Bool("add", false, "If adding new Todo")
	list := flag.Bool("list", false, "List all todos.")
	complete := flag.Int("complete", 0, "Todo Index (1 based) to be marked completed")

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoListFileName = os.Getenv("TODO_FILENAME")
	}

	todoList := &todoscli.TodoList{}
	loadErr := todoList.DeserializeAndLoad(todoListFileName)
	if loadErr != nil {
		fmt.Fprintln(os.Stderr, loadErr)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(todoList)
	case *complete > 0:
		// Mark the given todo task as complete.
		if err := todoList.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the modified list:
		if err := todoList.SerializeAndSave(todoListFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *addTodo:
		// Read and join to get the new todo task name.
		todoNameToAdd, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		todoList.Add(todoNameToAdd)
		// Save the new list
		if err := todoList.SerializeAndSave(todoListFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	// if no extra arg, then print the list.
	case len(os.Args) == 1:
		for _, todoItem := range *todoList {
			fmt.Println(todoItem.Task)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(scanner.Text()) == 0 {
		return "", fmt.Errorf("Todo cannot be blank")
	}
	return scanner.Text(), nil
}
