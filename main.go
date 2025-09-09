package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task represents a single task with all required properties
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TaskList represents the collection of tasks
type TaskList struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"nextId"`
}

const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
	TasksFile        = "tasks.json"
)

// loadTasks reads tasks from the JSON file
func loadTasks() (*TaskList, error) {
	taskList := &TaskList{
		Tasks:  []Task{},
		NextID: 1,
	}

	// Check if file exists
	if _, err := os.Stat(TasksFile); os.IsNotExist(err) {
		return taskList, nil
	}

	data, err := os.ReadFile(TasksFile)
	if err != nil {
		return nil, fmt.Errorf("error reading tasks file: %v", err)
	}

	if len(data) == 0 {
		return taskList, nil
	}

	err = json.Unmarshal(data, taskList)
	if err != nil {
		return nil, fmt.Errorf("error parsing tasks file: %v", err)
	}

	return taskList, nil
}

// saveTasks writes tasks to the JSON file
func saveTasks(taskList *TaskList) error {
	data, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling tasks: %v", err)
	}

	err = os.WriteFile(TasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing tasks file: %v", err)
	}

	return nil
}

// findTaskByID finds a task by its ID
func (tl *TaskList) findTaskByID(id int) *Task {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			return &tl.Tasks[i]
		}
	}
	return nil
}

// addTask adds a new task
func addTask(description string) error {
	taskList, err := loadTasks()
	if err != nil {
		return err
	}

	now := time.Now()
	task := Task{
		ID:          taskList.NextID,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskList.Tasks = append(taskList.Tasks, task)
	taskList.NextID++

	err = saveTasks(taskList)
	if err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
	return nil
}

// updateTask updates an existing task's description
func updateTask(id int, description string) error {
	taskList, err := loadTasks()
	if err != nil {
		return err
	}

	task := taskList.findTaskByID(id)
	if task == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	task.Description = description
	task.UpdatedAt = time.Now()

	err = saveTasks(taskList)
	if err != nil {
		return err
	}

	fmt.Printf("Task %d updated successfully\n", id)
	return nil
}

// deleteTask removes a task by ID
func deleteTask(id int) error {
	taskList, err := loadTasks()
	if err != nil {
		return err
	}

	found := false
	for i, task := range taskList.Tasks {
		if task.ID == id {
			taskList.Tasks = append(taskList.Tasks[:i], taskList.Tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	err = saveTasks(taskList)
	if err != nil {
		return err
	}

	fmt.Printf("Task %d deleted successfully\n", id)
	return nil
}

// markTask updates the status of a task
func markTask(id int, status string) error {
	taskList, err := loadTasks()
	if err != nil {
		return err
	}

	task := taskList.findTaskByID(id)
	if task == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	err = saveTasks(taskList)
	if err != nil {
		return err
	}

	fmt.Printf("Task %d marked as %s\n", id, status)
	return nil
}

// listTasks displays tasks based on the specified filter
func listTasks(filter string) error {
	taskList, err := loadTasks()
	if err != nil {
		return err
	}

	if len(taskList.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	filteredTasks := []Task{}
	for _, task := range taskList.Tasks {
		switch filter {
		case "":
			filteredTasks = append(filteredTasks, task)
		case "todo":
			if task.Status == StatusTodo {
				filteredTasks = append(filteredTasks, task)
			}
		case "in-progress":
			if task.Status == StatusInProgress {
				filteredTasks = append(filteredTasks, task)
			}
		case "done":
			if task.Status == StatusDone {
				filteredTasks = append(filteredTasks, task)
			}
		default:
			return fmt.Errorf("invalid filter: %s. Valid filters are: todo, in-progress, done", filter)
		}
	}

	if len(filteredTasks) == 0 {
		if filter != "" {
			fmt.Printf("No tasks found with status: %s\n", filter)
		} else {
			fmt.Println("No tasks found.")
		}
		return nil
	}

	// Print header
	fmt.Println("ID | Status      | Description")
	fmt.Println("---|-------------|------------")

	// Print tasks
	for _, task := range filteredTasks {
		status := task.Status
		if len(status) < 11 {
			status = status + strings.Repeat(" ", 11-len(status))
		}
		fmt.Printf("%-2d | %-11s | %s\n", task.ID, status, task.Description)
	}

	return nil
}

// printUsage displays the usage information
func printUsage() {
	fmt.Println("Task Tracker CLI")
	fmt.Println("Usage:")
	fmt.Println("  task-cli add \"description\"           - Add a new task")
	fmt.Println("  task-cli update <id> \"description\"   - Update task description")
	fmt.Println("  task-cli delete <id>                 - Delete a task")
	fmt.Println("  task-cli mark-in-progress <id>       - Mark task as in progress")
	fmt.Println("  task-cli mark-done <id>              - Mark task as done")
	fmt.Println("  task-cli list                        - List all tasks")
	fmt.Println("  task-cli list todo                   - List todo tasks")
	fmt.Println("  task-cli list in-progress            - List in-progress tasks")
	fmt.Println("  task-cli list done                   - List completed tasks")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Description is required for add command")
			printUsage()
			os.Exit(1)
		}
		err := addTask(os.Args[2])
		if err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			os.Exit(1)
		}

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Error: ID and description are required for update command")
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: Invalid task ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		err = updateTask(id, os.Args[3])
		if err != nil {
			fmt.Printf("Error updating task: %v\n", err)
			os.Exit(1)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: ID is required for delete command")
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: Invalid task ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		err = deleteTask(id)
		if err != nil {
			fmt.Printf("Error deleting task: %v\n", err)
			os.Exit(1)
		}

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Error: ID is required for mark-in-progress command")
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: Invalid task ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		err = markTask(id, StatusInProgress)
		if err != nil {
			fmt.Printf("Error marking task: %v\n", err)
			os.Exit(1)
		}

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Error: ID is required for mark-done command")
			printUsage()
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: Invalid task ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		err = markTask(id, StatusDone)
		if err != nil {
			fmt.Printf("Error marking task: %v\n", err)
			os.Exit(1)
		}

	case "list":
		filter := ""
		if len(os.Args) > 2 {
			filter = os.Args[2]
		}
		err := listTasks(filter)
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}
