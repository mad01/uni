package task

import (
	"os"
	"testing"
)

func TestTaskStore_AddTask(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "uni-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a task store with the temp directory
	store := &TaskStore{
		dataDir: tempDir,
		tasks:   []Task{},
	}

	// Test adding a task
	task, err := store.AddTask("Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected task ID 1, got %d", task.ID)
	}

	if task.Name != "Test Task" {
		t.Errorf("Expected task name 'Test Task', got '%s'", task.Name)
	}

	if task.Description != "Test Description" {
		t.Errorf("Expected task description 'Test Description', got '%s'", task.Description)
	}

	if task.Status != StatusOpen {
		t.Errorf("Expected task status 'open', got '%s'", task.Status)
	}

	// Test adding another task to verify ID increment
	task2, err := store.AddTask("Second Task", "")
	if err != nil {
		t.Fatalf("Failed to add second task: %v", err)
	}

	if task2.ID != 2 {
		t.Errorf("Expected second task ID 2, got %d", task2.ID)
	}
}

func TestTaskStore_GetTask(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "uni-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := &TaskStore{
		dataDir: tempDir,
		tasks:   []Task{},
	}

	// Add a task
	addedTask, err := store.AddTask("Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Get the task
	retrievedTask, err := store.GetTask(addedTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if retrievedTask.ID != addedTask.ID {
		t.Errorf("Expected task ID %d, got %d", addedTask.ID, retrievedTask.ID)
	}

	// Test getting non-existent task
	_, err = store.GetTask(999)
	if err == nil {
		t.Error("Expected error when getting non-existent task")
	}
}

func TestTaskStore_UpdateTaskStatus(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "uni-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := &TaskStore{
		dataDir: tempDir,
		tasks:   []Task{},
	}

	// Add a task
	task, err := store.AddTask("Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Update status to working
	updatedTask, err := store.UpdateTaskStatus(task.ID, StatusWorking)
	if err != nil {
		t.Fatalf("Failed to update task status: %v", err)
	}

	if updatedTask.Status != StatusWorking {
		t.Errorf("Expected status 'working', got '%s'", updatedTask.Status)
	}

	// Verify the task was actually updated in the store
	retrievedTask, err := store.GetTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to get updated task: %v", err)
	}

	if retrievedTask.Status != StatusWorking {
		t.Errorf("Expected status 'working' in store, got '%s'", retrievedTask.Status)
	}
}

func TestTaskStore_ListTasksWithFilter(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "uni-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store := &TaskStore{
		dataDir: tempDir,
		tasks:   []Task{},
	}

	// Add tasks with different statuses
	_, _ = store.AddTask("Open Task", "")
	task2, _ := store.AddTask("Working Task", "")
	task3, _ := store.AddTask("Done Task", "")
	task4, _ := store.AddTask("Blocked Task", "")

	store.UpdateTaskStatus(task2.ID, StatusWorking)
	store.UpdateTaskStatus(task3.ID, StatusDone)
	store.UpdateTaskStatus(task4.ID, StatusBlocked)

	// Test filtering left tasks
	leftTasks := store.ListTasksWithFilter(true, false)
	if len(leftTasks) != 3 {
		t.Errorf("Expected 3 left tasks, got %d", len(leftTasks))
	}

	// Test filtering closed tasks
	closedTasks := store.ListTasksWithFilter(false, true)
	if len(closedTasks) != 1 {
		t.Errorf("Expected 1 closed task, got %d", len(closedTasks))
	}

	// Test no filter
	allTasks := store.ListTasksWithFilter(false, false)
	if len(allTasks) != 4 {
		t.Errorf("Expected 4 total tasks, got %d", len(allTasks))
	}
}

func TestIsLeftStatus(t *testing.T) {
	leftStatuses := []TaskStatus{StatusOpen, StatusWorking, StatusBlocked}
	for _, status := range leftStatuses {
		if !isLeftStatus(status) {
			t.Errorf("Expected %s to be a left status", status)
		}
	}

	nonLeftStatuses := []TaskStatus{StatusDone, StatusCancel}
	for _, status := range nonLeftStatuses {
		if isLeftStatus(status) {
			t.Errorf("Expected %s to not be a left status", status)
		}
	}
}

func TestIsClosedStatus(t *testing.T) {
	closedStatuses := []TaskStatus{StatusDone, StatusCancel}
	for _, status := range closedStatuses {
		if !isClosedStatus(status) {
			t.Errorf("Expected %s to be a closed status", status)
		}
	}

	nonClosedStatuses := []TaskStatus{StatusOpen, StatusWorking, StatusBlocked}
	for _, status := range nonClosedStatuses {
		if isClosedStatus(status) {
			t.Errorf("Expected %s to not be a closed status", status)
		}
	}
}

func TestTaskStore_LoadAndSaveTasks(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "uni-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create first store and add tasks
	store1 := &TaskStore{
		dataDir: tempDir,
		tasks:   []Task{},
	}

	task1, err := store1.AddTask("Task 1", "Description 1")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	task2, err := store1.AddTask("Task 2", "Description 2")
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Create second store and load tasks
	store2 := &TaskStore{
		dataDir: tempDir,
		tasks:   []Task{},
	}

	err = store2.loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	if len(store2.tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(store2.tasks))
	}

	// Verify tasks were loaded correctly
	loadedTask1, err := store2.GetTask(task1.ID)
	if err != nil {
		t.Fatalf("Failed to get loaded task 1: %v", err)
	}

	if loadedTask1.Name != task1.Name {
		t.Errorf("Expected task name '%s', got '%s'", task1.Name, loadedTask1.Name)
	}

	loadedTask2, err := store2.GetTask(task2.ID)
	if err != nil {
		t.Fatalf("Failed to get loaded task 2: %v", err)
	}

	if loadedTask2.Name != task2.Name {
		t.Errorf("Expected task name '%s', got '%s'", task2.Name, loadedTask2.Name)
	}
}
