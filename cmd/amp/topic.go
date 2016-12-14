package main

import (
	"github.com/spf13/cobra"
)

// TopicCmd is the main command for attaching topic subcommands.
var TopicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Topic operations",
	Long:  `Manage topic-related operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return AMP.Connect()
	},
}

func init() {
	RootCmd.AddCommand(TopicCmd)
}
