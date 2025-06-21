package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mad01/uni/internal/task"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit <id>",
	Aliases: []string{"e"},
	Short:   "Edit a task using your default editor",
	Long:    `Edit a task by opening it in your default editor (set via EDITOR environment variable).`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		store, err := task.NewTaskStore()
		if err != nil {
			return err
		}

		taskToEdit, err := store.GetTask(id)
		if err != nil {
			return err
		}

		// Get editor from environment, default to vi
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		// Create a temporary file with current task content
		tempFile, err := os.CreateTemp("", fmt.Sprintf("uni-task-%d-*.txt", id))
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		// Write current task content to temp file
		content := fmt.Sprintf("Name: %s\nDescription: %s\n", taskToEdit.Name, taskToEdit.Description)
		if _, err := tempFile.WriteString(content); err != nil {
			return fmt.Errorf("failed to write to temporary file: %v", err)
		}
		tempFile.Close()

		// Store original modification time
		stat, err := os.Stat(tempFile.Name())
		if err != nil {
			return fmt.Errorf("failed to get file stats: %v", err)
		}
		originalModTime := stat.ModTime()

		// Open editor
		editorCmd := exec.Command(editor, tempFile.Name())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("editor failed: %v", err)
		}

		// Check if file was modified
		newStat, err := os.Stat(tempFile.Name())
		if err != nil {
			return fmt.Errorf("failed to get file stats after edit: %v", err)
		}

		if newStat.ModTime().Equal(originalModTime) {
			fmt.Println("No changes made.")
			return nil
		}

		// Read the edited content
		editedContent, err := os.ReadFile(tempFile.Name())
		if err != nil {
			return fmt.Errorf("failed to read edited file: %v", err)
		}

		// Parse the edited content
		newName, newDescription, err := parseEditedContent(string(editedContent))
		if err != nil {
			return err
		}

		// Update the task
		taskToEdit.Name = newName
		taskToEdit.Description = newDescription
		taskToEdit.UpdatedAt = time.Now()

		// Update task in store
		if err := store.UpdateTask(taskToEdit); err != nil {
			return fmt.Errorf("failed to update task: %v", err)
		}

		fmt.Printf("Task #%d updated successfully.\n", taskToEdit.ID)
		return nil
	},
}

func parseEditedContent(content string) (name, description string, err error) {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name: ") {
			name = strings.TrimPrefix(line, "Name: ")
		} else if strings.HasPrefix(line, "Description: ") {
			description = strings.TrimPrefix(line, "Description: ")
		}
	}

	if name == "" {
		return "", "", fmt.Errorf("task name cannot be empty")
	}

	return name, description, nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}
