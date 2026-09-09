package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do [task number]",
	Short: "Mark a task on your TODO list as complete",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// We'll parse the task ID here once we hook up BoltDB
		fmt.Printf("You have completed task #%s.\n", args[0])
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}