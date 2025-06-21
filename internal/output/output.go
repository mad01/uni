package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mad01/uni/internal/task"
	"gopkg.in/yaml.v3"
)

// FormatTasks formats tasks according to the specified output format
func FormatTasks(tasks []task.Task, format string) error {
	switch format {
	case "json":
		return formatJSON(tasks)
	case "yaml":
		return formatYAML(tasks)
	case "text":
		return formatText(tasks)
	case "normal":
		return formatNormal(tasks)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// FormatTask formats a single task according to the specified output format
func FormatTask(t *task.Task, format string) error {
	switch format {
	case "json":
		return formatJSON([]*task.Task{t})
	case "yaml":
		return formatYAML([]*task.Task{t})
	case "text":
		return formatText([]*task.Task{t})
	case "normal":
		return formatNormal([]*task.Task{t})
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func formatJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func formatYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(data)
}

func formatText(data interface{}) error {
	switch v := data.(type) {
	case []task.Task:
		return formatTasksText(v)
	case []*task.Task:
		tasks := make([]task.Task, len(v))
		for i, t := range v {
			tasks[i] = *t
		}
		return formatTasksText(tasks)
	default:
		return fmt.Errorf("unsupported data type for text format")
	}
}

func formatNormal(data interface{}) error {
	switch v := data.(type) {
	case []task.Task:
		return formatTasksNormal(v)
	case []*task.Task:
		tasks := make([]task.Task, len(v))
		for i, t := range v {
			tasks[i] = *t
		}
		return formatTasksNormal(tasks)
	default:
		return fmt.Errorf("unsupported data type for normal format")
	}
}

func formatTasksText(tasks []task.Task) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tNAME\tDESCRIPTION")
	for _, t := range tasks {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, strings.ToUpper(string(t.Status)), t.Name, t.Description)
	}
	return w.Flush()
}

func formatTasksNormal(tasks []task.Task) error {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	for _, t := range tasks {
		statusColor := getStatusColor(t.Status)
		fmt.Printf("%s#%d%s [%s%s%s] %s\n",
			"\033[1m", t.ID, "\033[0m",
			statusColor, strings.ToUpper(string(t.Status)), "\033[0m",
			t.Name)
		if t.Description != "" {
			fmt.Printf("  %s\n", t.Description)
		}
		fmt.Println()
	}
	return nil
}

func getStatusColor(status task.TaskStatus) string {
	switch status {
	case task.StatusOpen:
		return "\033[33m" // Yellow
	case task.StatusWorking:
		return "\033[34m" // Blue
	case task.StatusBlocked:
		return "\033[35m" // Magenta
	case task.StatusDone:
		return "\033[32m" // Green
	case task.StatusCancel:
		return "\033[31m" // Red
	default:
		return "\033[0m" // Reset
	}
}
