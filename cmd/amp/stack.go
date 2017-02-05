package main

import (
	"github.com/spf13/cobra"
)

// StackCmd is the main command for attaching stack subcommands.
var StackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Stack operations",
	Long:  `Stack command manages all stack-related operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(StackCmd)
}
