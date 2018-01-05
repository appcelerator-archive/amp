package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewNodeCommand returns a new instance of the cluster node command.
func NewNodeCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node",
		Short:   "Node management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewNodeListCommand(c))
	cmd.AddCommand(NewNodeInspectCommand(c))
	cmd.AddCommand(NewNodeCleanupCommand(c))
	return cmd
}
