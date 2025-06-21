package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
	showLeft     bool
	showClosed   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "uni",
	Short: "A minimal task management CLI",
	Long:  `uni is a minimal task management CLI that stores tasks in JSON files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "normal", "Output format (normal, text, json, yaml)")
	rootCmd.PersistentFlags().BoolVar(&showLeft, "left", false, "Show only left (open, working, blocked) tasks")
	rootCmd.PersistentFlags().BoolVar(&showClosed, "closed", false, "Show only closed (done, cancelled) tasks")
}

// GetOutputFormat returns the current output format
func GetOutputFormat() string {
	return outputFormat
}

// GetShowLeft returns if only left tasks should be shown
func GetShowLeft() bool {
	return showLeft
}

// GetShowClosed returns if only closed tasks should be shown
func GetShowClosed() bool {
	return showClosed
}

// ValidateOutputFormat validates the output format
func ValidateOutputFormat(format string) error {
	validFormats := []string{"normal", "text", "json", "yaml"}
	for _, valid := range validFormats {
		if format == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid output format: %s. Valid formats: %v", format, validFormats)
}
