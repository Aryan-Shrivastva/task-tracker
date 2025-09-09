# Task Tracker CLI

A simple command-line interface (CLI) application built in Go for tracking and managing your tasks. This project helps you keep track of what you need to do, what you have done, and what you are currently working on.

## Features

- ✅ Add new tasks
- ✅ Update task descriptions
- ✅ Delete tasks
- ✅ Mark tasks as in-progress or done
- ✅ List all tasks
- ✅ Filter tasks by status (todo, in-progress, done)
- ✅ Persistent storage using JSON file
- ✅ Automatic timestamp tracking (created and updated times)
- ✅ Auto-incrementing task IDs
- ✅ Comprehensive error handling

## Installation

### Prerequisites

- Go 1.16 or higher installed on your system

### Building the Application

1. Clone or download this repository
2. Navigate to the project directory
3. Build the application:

```bash
go build -o task-cli main.go
```

On Windows, the executable will be named `task-cli.exe`.

## Usage

The application accepts commands and arguments from the command line. All tasks are stored in a `tasks.json` file in the current directory.

### Commands

#### Add a new task
```bash
./task-cli add "Buy groceries"
# Output: Task added successfully (ID: 1)
```

#### List all tasks
```bash
./task-cli list
```

#### List tasks by status
```bash
./task-cli list todo          # List todo tasks
./task-cli list in-progress   # List in-progress tasks
./task-cli list done          # List completed tasks
```

#### Update a task description
```bash
./task-cli update 1 "Buy groceries and cook dinner"
# Output: Task 1 updated successfully
```

#### Mark task status
```bash
./task-cli mark-in-progress 1    # Mark task as in-progress
./task-cli mark-done 1           # Mark task as done
```

#### Delete a task
```bash
./task-cli delete 1
# Output: Task 1 deleted successfully
```

#### Help
```bash
./task-cli
# Shows usage information and available commands
```

## Task Properties

Each task contains the following properties:

- **id**: Unique identifier for the task (auto-generated)
- **description**: Short description of the task
- **status**: Current status (`todo`, `in-progress`, `done`)
- **createdAt**: Timestamp when the task was created
- **updatedAt**: Timestamp when the task was last modified

## File Structure

```
task-tracker/
├── main.go           # Main application code
├── go.mod           # Go module file
├── task-cli.exe     # Compiled executable (after building)
├── tasks.json       # JSON file storing all tasks (auto-created)
└── README.md        # This file
```

## Example Usage Session

```bash
# Add some tasks
./task-cli add "Buy groceries"
./task-cli add "Complete Go project"
./task-cli add "Write documentation"

# List all tasks
./task-cli list
# ID | Status      | Description
# ---|-------------|------------
# 1  | todo        | Buy groceries
# 2  | todo        | Complete Go project
# 3  | todo        | Write documentation

# Mark a task as in-progress
./task-cli mark-in-progress 2

# Mark a task as done
./task-cli mark-done 1

# Update a task description
./task-cli update 3 "Write comprehensive documentation with examples"

# List tasks by status
./task-cli list done
./task-cli list in-progress
./task-cli list todo

# Delete a completed task
./task-cli delete 1
```

## JSON File Format

Tasks are stored in `tasks.json` with the following structure:

```json
{
  "tasks": [
    {
      "id": 1,
      "description": "Buy groceries",
      "status": "todo",
      "createdAt": "2025-09-10T02:08:50.123456+05:30",
      "updatedAt": "2025-09-10T02:08:50.123456+05:30"
    }
  ],
  "nextId": 2
}
```

## Error Handling

The application handles various error scenarios:

- Invalid command arguments
- Non-existent task IDs
- Invalid task IDs (non-numeric)
- File system errors
- JSON parsing errors
- Missing required arguments

## Development

### Project Structure

The application is built as a single Go file (`main.go`) with the following key components:

- **Task struct**: Defines the task data structure
- **TaskList struct**: Manages the collection of tasks and ID generation
- **File operations**: Functions for reading and writing the JSON file
- **Command handlers**: Functions for each CLI command
- **Main function**: Command-line argument parsing and routing

### Key Functions

- `loadTasks()`: Reads tasks from the JSON file
- `saveTasks()`: Writes tasks to the JSON file
- `addTask()`: Adds a new task
- `updateTask()`: Updates an existing task
- `deleteTask()`: Removes a task
- `markTask()`: Changes task status
- `listTasks()`: Displays tasks with optional filtering

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is open source and available under the MIT License.

## Requirements Met

✅ **Command Line Interface**: Accepts user actions and inputs as arguments  
✅ **JSON Storage**: Stores tasks in a JSON file in the current directory  
✅ **File Creation**: Creates JSON file if it doesn't exist  
✅ **Native File System**: Uses Go's native `os` package for file operations  
✅ **No External Dependencies**: Built using only Go standard library  
✅ **Error Handling**: Comprehensive error handling for all operations  
✅ **All Required Features**:
- Add, update, and delete tasks
- Mark tasks as in-progress or done
- List all tasks
- List tasks by status (done, todo, in-progress)
- Proper task properties (id, description, status, createdAt, updatedAt)

## Future Enhancements

- Task due dates and reminders
- Task priorities
- Task categories/tags
- Search functionality
- Export options (CSV, plain text)
- Task completion statistics
- Colored output for better readability