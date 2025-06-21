package cmd

import (
	"github.com/mad01/uni/internal/output"
	"github.com/mad01/uni/internal/task"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all tasks",
	Long:    `List all tasks with their current status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateOutputFormat(GetOutputFormat()); err != nil {
			return err
		}

		store, err := task.NewTaskStore()
		if err != nil {
			return err
		}

		tasks := store.ListTasksWithFilter(GetShowLeft(), GetShowClosed())
		return output.FormatTasks(tasks, GetOutputFormat())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
