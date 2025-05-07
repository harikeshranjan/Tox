package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reindexCmd represents the reindex command
var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex all todos",
	Long:  `Reset and reorder all todo IDs to be sequential starting from 1.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := dbManager.ReindexAll()
		if err != nil {
			return err
		}

		fmt.Println("Successfully reindexed all todos")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reindexCmd)
}
