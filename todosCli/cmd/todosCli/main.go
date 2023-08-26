package main

import (
	"flag"
	"fmt"
	"os"
	todoscli "piyush-mishra-pm/go-cli-tools/todosCli"
)

// Hardcoded save file name:
const todoListFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed by Piyush Mishra\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parse Flags:
	addTodoTaskName := flag.String("add", "", "Todo item to be added in TodoList")
	list := flag.Bool("list", false, "List all todos.")
	complete := flag.Int("complete", 0, "Todo Index (1 based) to be marked completed")

	flag.Parse()

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
	case *addTodoTaskName != "":
		todoList.Add(*addTodoTaskName)
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
