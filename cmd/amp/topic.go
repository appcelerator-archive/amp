package main

import (
	"github.com/spf13/cobra"
)

// TopicCmd is the main command for attaching topic subcommands.
var TopicCmd = &cobra.Command{
	Use:   "topic operations",
	Short: "topic operations",
	Long:  `Manage topic-related operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = AMP.Connect()
		return
	},
}

func init() {
	RootCmd.AddCommand(TopicCmd)
}
