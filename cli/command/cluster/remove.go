package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

var (
	destroyArgs = []string{"-d"}
)

// NewRemoveCommand returns a new instance of the remove command for destroying a cluster.
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "destroy",
		Short:   "Destroy an amp cluster",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroy(c)
		},
	}
}

func destroy(c cli.Interface) error {
	return updateCluster(c, destroyArgs)
}
