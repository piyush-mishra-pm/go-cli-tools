package main

import (
	"fmt"
	"os"
	todoscli "piyush-mishra-pm/go-cli-tools/todosCli"
	"strings"
)

// Hardcoded save file name:
const todoListFileName = ".todo.json"

func main() {
	todoList := &todoscli.TodoList{}
	loadErr := todoList.DeserializeAndLoad(todoListFileName)
	if loadErr != nil {
		fmt.Fprintln(os.Stderr, loadErr)
		os.Exit(1)
	}

	switch {
	// if no extra arg, then print the list.
	case len(os.Args) == 1:
		for _, todoItem := range *todoList {
			fmt.Println(todoItem.Task)
		}
	default:
		// Concat all arg with space and add it to back of todo list.
		newTodoTaskName := strings.Join(os.Args[1:], " ")
		todoList.Add(newTodoTaskName)

		if saveErr := todoList.SerializeAndSave(todoListFileName); saveErr != nil {
			fmt.Fprintln(os.Stderr, saveErr)
			os.Exit(1)
		}
	}
}
