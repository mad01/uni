package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	StatusOpen    TaskStatus = "open"
	StatusDone    TaskStatus = "done"
	StatusBlocked TaskStatus = "blocked"
	StatusWorking TaskStatus = "working"
	StatusCancel  TaskStatus = "cancel"
)

// Task represents a single task
type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaskStore manages tasks
type TaskStore struct {
	dataDir string
	tasks   []Task
}

// NewTaskStore creates a new task store
func NewTaskStore() (*TaskStore, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	store := &TaskStore{
		dataDir: dataDir,
		tasks:   []Task{},
	}

	if err := store.ensureDataDir(); err != nil {
		return nil, err
	}

	if err := store.loadTasks(); err != nil {
		return nil, err
	}

	return store, nil
}

// getDataDir determines the data directory (.uni in git repo or ~/.uni)
func getDataDir() (string, error) {
	// Check if we're in a git repository and have a .uni directory
	if _, err := os.Stat(".git"); err == nil {
		uniDir := ".uni"
		if _, err := os.Stat(uniDir); err == nil {
			return uniDir, nil
		}
	}

	// Use ~/.uni as default
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".uni"), nil
}

// ensureDataDir creates the data directory if it doesn't exist
func (ts *TaskStore) ensureDataDir() error {
	return os.MkdirAll(ts.dataDir, 0755)
}

// getTasksFile returns the path to the tasks.json file
func (ts *TaskStore) getTasksFile() string {
	return filepath.Join(ts.dataDir, "tasks.json")
}

// loadTasks loads tasks from the JSON file
func (ts *TaskStore) loadTasks() error {
	tasksFile := ts.getTasksFile()

	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		// File doesn't exist, start with empty slice
		ts.tasks = []Task{}
		return nil
	}

	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		ts.tasks = []Task{}
		return nil
	}

	return json.Unmarshal(data, &ts.tasks)
}

// saveTasks saves tasks to the JSON file
func (ts *TaskStore) saveTasks() error {
	data, err := json.MarshalIndent(ts.tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ts.getTasksFile(), data, 0644)
}

// getNextID returns the next available ID
func (ts *TaskStore) getNextID() int {
	maxID := 0
	for _, task := range ts.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// AddTask adds a new task
func (ts *TaskStore) AddTask(name, description string) (*Task, error) {
	now := time.Now()
	task := Task{
		ID:          ts.getNextID(),
		Name:        name,
		Description: description,
		Status:      StatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ts.tasks = append(ts.tasks, task)

	if err := ts.saveTasks(); err != nil {
		return nil, err
	}

	return &task, nil
}

// GetTask gets a task by ID
func (ts *TaskStore) GetTask(id int) (*Task, error) {
	for i, task := range ts.tasks {
		if task.ID == id {
			return &ts.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// ListTasks returns all tasks sorted by ID
func (ts *TaskStore) ListTasks() []Task {
	return ts.ListTasksWithFilter(false, false)
}

// ListTasksWithFilter returns tasks with optional filtering
func (ts *TaskStore) ListTasksWithFilter(showLeft, showClosed bool) []Task {
	// Create a copy and sort by ID
	tasks := make([]Task, len(ts.tasks))
	copy(tasks, ts.tasks)

	// Apply filters
	if showLeft || showClosed {
		filtered := []Task{}
		for _, task := range tasks {
			if showLeft && isLeftStatus(task.Status) {
				filtered = append(filtered, task)
			} else if showClosed && isClosedStatus(task.Status) {
				filtered = append(filtered, task)
			}
		}
		tasks = filtered
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks
}

// isLeftStatus returns true if the status represents a "left" (open) task
func isLeftStatus(status TaskStatus) bool {
	return status == StatusOpen || status == StatusWorking || status == StatusBlocked
}

// isClosedStatus returns true if the status represents a closed task
func isClosedStatus(status TaskStatus) bool {
	return status == StatusDone || status == StatusCancel
}

// UpdateTaskStatus updates the status of a task
func (ts *TaskStore) UpdateTaskStatus(id int, status TaskStatus) (*Task, error) {
	for i, task := range ts.tasks {
		if task.ID == id {
			ts.tasks[i].Status = status
			ts.tasks[i].UpdatedAt = time.Now()

			if err := ts.saveTasks(); err != nil {
				return nil, err
			}

			return &ts.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// UpdateTask updates a task's name and description
func (ts *TaskStore) UpdateTask(updatedTask *Task) error {
	for i, task := range ts.tasks {
		if task.ID == updatedTask.ID {
			ts.tasks[i] = *updatedTask
			return ts.saveTasks()
		}
	}
	return fmt.Errorf("task with ID %d not found", updatedTask.ID)
}
