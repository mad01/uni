package cmd

import (
	"fmt"

	"github.com/mad01/uni/internal/output"
	"github.com/mad01/uni/internal/task"
	"github.com/spf13/cobra"
)

var (
	addName        string
	addDescription string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add a new task",
	Long:    `Add a new task with a name and optional description using flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateOutputFormat(GetOutputFormat()); err != nil {
			return err
		}

		if addName == "" {
			return fmt.Errorf("task name is required (use --name or -n)")
		}

		store, err := task.NewTaskStore()
		if err != nil {
			return err
		}

		newTask, err := store.AddTask(addName, addDescription)
		if err != nil {
			return err
		}

		if GetOutputFormat() == "normal" {
			fmt.Printf("Task #%d created successfully.\n", newTask.ID)
			return nil
		}

		return output.FormatTask(newTask, GetOutputFormat())
	},
}

func init() {
	addCmd.Flags().StringVarP(&addName, "name", "n", "", "Task name (required)")
	addCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Task description (optional)")
	addCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(addCmd)
}
