package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewRemoveCommand returns a new instance of the remove command for destroying a cluster.
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm",
		Aliases: []string{"remove", "destroy"},
		Short:   "Destroy an amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c)
		},
	}
}

func remove(c cli.Interface) error {
	// TODO: only supporting local cluster management for this release
	args := []string{"bootstrap/bootstrap", "-d", DefaultLocalClusterID}
	return queryCluster(c, args)
}
