package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You have the following tasks:")
		fmt.Println("1. review talk proposal")
		fmt.Println("2. clean dishes")
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}