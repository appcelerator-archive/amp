package main

import (
	"github.com/spf13/cobra"
)

// ClusterCmd is the main command for attaching cluster subcommands.
var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster operations",
	Long:  `Manage cluster-related operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ClusterCmd)
}
