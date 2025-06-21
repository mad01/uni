package cmd

import (
	"fmt"
	"strconv"

	"github.com/mad01/uni/internal/output"
	"github.com/mad01/uni/internal/task"
	"github.com/spf13/cobra"
)

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:     "done <id>",
	Aliases: []string{"d"},
	Short:   "Mark a task as done",
	Long:    `Mark a task as done by providing its ID.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateOutputFormat(GetOutputFormat()); err != nil {
			return err
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		store, err := task.NewTaskStore()
		if err != nil {
			return err
		}

		updatedTask, err := store.UpdateTaskStatus(id, task.StatusDone)
		if err != nil {
			return err
		}

		if GetOutputFormat() == "normal" {
			fmt.Printf("Task #%d marked as done.\n", updatedTask.ID)
			return nil
		}

		return output.FormatTask(updatedTask, GetOutputFormat())
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
