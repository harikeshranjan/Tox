/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your todos",
	Long:  `Display a list of your todos, by default only showing incomplete tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showDone, _ := cmd.Flags().GetBool("all")

		todos, err := dbManager.ListTodos(showDone)
		if err != nil {
			return err
		}

		if len(todos) == 0 {
			fmt.Println("No todos found!")
			return nil
		}

		fmt.Println("ID  STATUS  TASK")
		fmt.Println("--  ------  ----")

		for _, todo := range todos {
			status := "[ ]"
			if todo.Done {
				status = "[✓]"
			}

			fmt.Printf("%2d  %s  %s", todo.ID, status, todo.Task)

			if todo.Done && todo.DoneAt != nil {
				fmt.Printf(" (completed %s)", formatTimeAgo(*todo.DoneAt))
			}

			fmt.Println()
		}

		return nil
	},
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "Show all todos, including completed ones")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
