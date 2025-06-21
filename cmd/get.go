package cmd

import (
	"fmt"
	"strconv"

	"github.com/mad01/uni/internal/output"
	"github.com/mad01/uni/internal/task"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a specific task",
	Long:  `Get a specific task by providing its ID.`,
	Args:  cobra.ExactArgs(1),
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

		t, err := store.GetTask(id)
		if err != nil {
			return err
		}

		return output.FormatTask(t, GetOutputFormat())
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
