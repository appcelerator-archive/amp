package main

import (
	"github.com/spf13/cobra"
)

// NodeCmd is the main command for attaching cluster subcommands.
var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node operations",
	Long:  `Manage node-related operations.`,
	//Aliases: []string{"pf"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return client.initConnection()
	},
}

func init() {
	RootCmd.AddCommand(NodeCmd)
}
