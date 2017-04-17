package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewNodeCommand returns a new instance of the cluster node command.
func NewNodeCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node",
		Short:   "Cluster node management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewNodeListCommand(c))
	return cmd
}
