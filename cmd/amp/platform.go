package main

import (
	"github.com/spf13/cobra"
)

// PlatformCmd is the main command for attaching topic subcommands.
var PlatformCmd = &cobra.Command{
	Use:     "platform",
	Short:   "Platform operations (alias: pf)",
	Long:    `Manage platform-related operations.`,
	Aliases: []string{"pf"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(PlatformCmd)
}
