package todoscli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// todo item.
type todo struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// Todo list.
type TodoList []todo

// Add todo to TodoList
func (tl *TodoList) Add(task string) {
	newTask := todo{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*tl = append(*tl, newTask)
}

// Mark todo as Complete:
// Accepts 1 based index (Not 0 based).
func (tl *TodoList) Complete(todoIndex int) error {
	if isIndexOutOfRange(todoIndex, tl) {
		return fmt.Errorf("Todo item %d doesnt exist in TodoList. Cant mark Complete", todoIndex)
	}

	(*tl)[todoIndex-1].Done = true
	(*tl)[todoIndex-1].CompletedAt = time.Now()

	return nil
}

func (tl *TodoList) Delete(todoIndex int) error {
	if isIndexOutOfRange(todoIndex, tl) {
		return fmt.Errorf("Todo item %d doesnt exist in TodoList. Cant Delete.", todoIndex)
	}

	*tl = append((*tl)[:todoIndex-1], (*tl)[todoIndex:]...)
	return nil
}

func (tl *TodoList) SerializeAndSave(filename string) error {
	// Serialise as JSON:
	jsonTodoList, err := json.Marshal(tl)
	if err != nil {
		return err
	}

	// Save to file:
	return os.WriteFile(filename, jsonTodoList, 0644)
}

func (tl *TodoList) DeserializeAndLoad(filename string) error {
	// Load file
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		// File Doesn't exist:
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	// File exists, but is empty:
	if len(fileContent) == 0 {
		return nil
	}

	// Deserialize JSON to list.
	return json.Unmarshal(fileContent, tl)
}

// Checks if 1 based index is out of range.
func isIndexOutOfRange(index int, tl *TodoList) bool {
	if index <= 0 || index > len(*tl) {
		return true
	}
	return false
}
