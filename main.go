package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
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

// Global taskList to avoid repeated load/save operations
var taskList *TaskList

// fatal prints an error message and exits with code 1
func fatal(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
	} else {
		fmt.Println(msg)
	}
	os.Exit(1)
}

// parseID parses a string argument to an integer ID
func parseID(arg string) (int, error) {
	return strconv.Atoi(arg)
}

// requireArgs checks if minimum number of arguments are provided
func requireArgs(min int, usage string) {
	if len(os.Args) < min {
		fatal(fmt.Sprintf("Error: %s", usage), nil)
	}
}

// loadTasks reads tasks from the JSON file
func loadTasks() (*TaskList, error) {
	tl := &TaskList{
		Tasks:  []Task{},
		NextID: 1,
	}

	// Check if file exists
	if _, err := os.Stat(TasksFile); os.IsNotExist(err) {
		return tl, nil
	}

	data, err := os.ReadFile(TasksFile)
	if err != nil {
		return nil, fmt.Errorf("error reading tasks file: %v", err)
	}

	if len(data) == 0 {
		return tl, nil
	}

	err = json.Unmarshal(data, tl)
	if err != nil {
		return nil, fmt.Errorf("error parsing tasks file: %v", err)
	}

	return tl, nil
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

	fmt.Printf("Task added successfully (ID: %d)\n", task.ID)
	return nil
}

// updateTask updates an existing task's description
func updateTask(id int, description string) error {
	task := taskList.findTaskByID(id)
	if task == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	task.Description = description
	task.UpdatedAt = time.Now()

	fmt.Printf("Task %d updated successfully\n", id)
	return nil
}

// deleteTask removes a task by ID
func deleteTask(id int) error {
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

	fmt.Printf("Task %d deleted successfully\n", id)
	return nil
}

// markTask updates the status of a task
func markTask(id int, status string) error {
	task := taskList.findTaskByID(id)
	if task == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	fmt.Printf("Task %d marked as %s\n", id, status)
	return nil
}

// listTasks displays tasks based on the specified filter
func listTasks(filter string) error {
	if len(taskList.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	// Define filter functions
	filters := map[string]func(Task) bool{
		"":               func(t Task) bool { return true },
		StatusTodo:       func(t Task) bool { return t.Status == StatusTodo },
		StatusInProgress: func(t Task) bool { return t.Status == StatusInProgress },
		StatusDone:       func(t Task) bool { return t.Status == StatusDone },
	}

	match, ok := filters[filter]
	if !ok {
		return fmt.Errorf("invalid filter: %s. Valid filters are: %s, %s, %s", filter, StatusTodo, StatusInProgress, StatusDone)
	}

	// Filter and count tasks
	filteredTasks := []Task{}
	for _, task := range taskList.Tasks {
		if match(task) {
			filteredTasks = append(filteredTasks, task)
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

	// Use tabwriter for pretty output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tDescription")
	fmt.Fprintln(w, "---\t------\t-----------")
	for _, task := range filteredTasks {
		fmt.Fprintf(w, "%d\t%s\t%s\n", task.ID, task.Status, task.Description)
	}
	w.Flush()

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
	fmt.Printf("  task-cli list %s                   - List %s tasks\n", StatusTodo, StatusTodo)
	fmt.Printf("  task-cli list %s            - List %s tasks\n", StatusInProgress, StatusInProgress)
	fmt.Printf("  task-cli list %s                   - List %s tasks\n", StatusDone, StatusDone)
}

func main() {
	// Load tasks once at startup
	var err error
	taskList, err = loadTasks()
	if err != nil {
		fatal("Error loading tasks", err)
	}

	// Ensure we save tasks before exit
	defer func() {
		if err := saveTasks(taskList); err != nil {
			fmt.Printf("Warning: Error saving tasks: %v\n", err)
		}
	}()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		requireArgs(3, "Description is required for add command")
		if err := addTask(os.Args[2]); err != nil {
			fatal("Error adding task", err)
		}

	case "update":
		requireArgs(4, "ID and description are required for update command")
		id, err := parseID(os.Args[2])
		if err != nil {
			fatal(fmt.Sprintf("Invalid task ID: %s", os.Args[2]), err)
		}
		if err := updateTask(id, os.Args[3]); err != nil {
			fatal("Error updating task", err)
		}

	case "delete":
		requireArgs(3, "ID is required for delete command")
		id, err := parseID(os.Args[2])
		if err != nil {
			fatal(fmt.Sprintf("Invalid task ID: %s", os.Args[2]), err)
		}
		if err := deleteTask(id); err != nil {
			fatal("Error deleting task", err)
		}

	case "mark-in-progress":
		requireArgs(3, "ID is required for mark-in-progress command")
		id, err := parseID(os.Args[2])
		if err != nil {
			fatal(fmt.Sprintf("Invalid task ID: %s", os.Args[2]), err)
		}
		if err := markTask(id, StatusInProgress); err != nil {
			fatal("Error marking task", err)
		}

	case "mark-done":
		requireArgs(3, "ID is required for mark-done command")
		id, err := parseID(os.Args[2])
		if err != nil {
			fatal(fmt.Sprintf("Invalid task ID: %s", os.Args[2]), err)
		}
		if err := markTask(id, StatusDone); err != nil {
			fatal("Error marking task", err)
		}

	case "list":
		filter := ""
		if len(os.Args) > 2 {
			filter = os.Args[2]
		}
		if err := listTasks(filter); err != nil {
			fatal("Error listing tasks", err)
		}

	default:
		fatal(fmt.Sprintf("Unknown command '%s'", command), nil)
		printUsage()
	}
}
