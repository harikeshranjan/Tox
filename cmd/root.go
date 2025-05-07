/*
Copyright © 2025 Harikesh Ranjan Sinha <ranjansinhaharikesh@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/harikeshranjan/tox/db"
	"github.com/spf13/cobra"
)

var dbManager *db.Manager

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tox",
	Short: "A simple CLI Todo Manager",
	Long: `Tox is a CLI tool to manage your todos from the terminal.
You can add, list, mark as done, and delete todos with simple commands.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		dbManager, err = db.NewManager()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if dbManager != nil {
			if err := dbManager.Close(); err != nil {
				fmt.Println("Error closing the database:", err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tox.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
