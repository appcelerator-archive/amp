package main

import (
	"github.com/spf13/cobra"
)

// StackCmd is the main command for attaching stack subcommands.
var StackCmd = &cobra.Command{
	Use:   "stack operations",
	Short: "stack operations",
	Long:  `Manage stack-related operations.`,
}

func init() {
	RootCmd.AddCommand(StackCmd)
}
