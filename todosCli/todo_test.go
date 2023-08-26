package todoscli_test

import (
	"os"
	todoscli "piyush-mishra-pm/go-cli-tools/todosCli"
	"testing"
)

func TestAdd(t *testing.T) {
	todoList := todoscli.TodoList{}

	todoTaskName := "New Todo Task"
	todoList.Add(todoTaskName)

	if len(todoList) != 1 {
		t.Errorf("Expected todolist length 1, got %q instead.", len(todoList))
	}
	if todoList[0].Task != todoTaskName {
		t.Errorf("Expected task name %q, got %q instead.", todoTaskName, todoList[0].Task)
	}
}

func TestComplete(t *testing.T) {
	todoList := todoscli.TodoList{}

	todoTaskName := "New Todo Task"
	todoList.Add(todoTaskName)

	if len(todoList) != 1 {
		t.Errorf("Expected todolist length 1, got %q instead.", len(todoList))
	}
	if todoList[0].Task != todoTaskName {
		t.Errorf("Expected task name %q, got %q instead.", todoTaskName, todoList[0].Task)
	}

	if todoList[0].Done {
		t.Errorf("Expected new task to be not completed")
	}
	err := todoList.Complete(1)
	if err != nil {
		t.Errorf("Error occurred during marking task as complete. Is index out of range?")
	}
	if !todoList[0].Done {
		t.Errorf("Expected task to be completed, after running Complete on it.")
	}
}

func TestDelete(t *testing.T) {
	todoList := todoscli.TodoList{}

	todoTaskNames := []string{"New Todo Task 1", "New Todo Task 2", "New Todo Task 3"}

	for _, taskName := range todoTaskNames {
		todoList.Add(taskName)
	}

	if len(todoList) != 3 {
		t.Errorf("Expected todolist length 3, got %q instead.", len(todoList))
	}
	if todoList[0].Task != todoTaskNames[0] {
		t.Errorf("Expected task name %q, got %q instead.", todoTaskNames[0], todoList[0].Task)
	}
	err := todoList.Delete(2)
	if err != nil {
		t.Errorf("Error occurred during deleting todo from list: %s. \nIs index out of range?", err)
	}
	if len(todoList) != 2 {
		t.Errorf("After deleting 1 element, expected todolist length 2, got %q instead.", len(todoList))
	}

	if todoList[1].Task != todoTaskNames[2] {
		t.Errorf("After deleting 2nd todo, expected 2nd todo name %q, but got %q instead", todoTaskNames[2], todoList[1].Task)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Testing Serialising and Saving:
	todoList1 := todoscli.TodoList{}

	todoTaskName := "New Todo Task Name"
	todoList1.Add(todoTaskName)

	if todoList1[0].Task != todoTaskName {
		t.Errorf("Expected task name %q, got %q instead", todoTaskName, todoList1[0].Task)
	}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Cant test Save & Load. Error creating temp file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	if err := todoList1.SerializeAndSave(tempFile.Name()); err != nil {
		t.Fatalf("Error serialising and saving todolist to file: %s", err)
	}

	// Testing Loading and Deserielising:
	todoList2 := todoscli.TodoList{}
	if err := todoList2.DeserializeAndLoad(tempFile.Name()); err != nil {
		t.Errorf("Error loading todo list from file: %s", err)
	}
	if todoList1[0].Task != todoList2[0].Task {
		t.Errorf("Saved task name didnt match with Loaded task name. Saved: %q, Loaded: %q", todoList1[0].Task, todoList2[0].Task)
	}
}
