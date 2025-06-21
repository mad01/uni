package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build-time variables (injected via ldflags)
var (
	gitHash = "dev"
	dirty   = "false"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show git version information",
	Long:    `Show the git commit hash that was used to build this binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("uni %s", gitHash)
		if dirty == "true" {
			fmt.Printf("-dirty")
		}
		fmt.Println()

		if date != "unknown" {
			fmt.Printf("built: %s\n", date)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
