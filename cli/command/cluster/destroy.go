package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

var (
	destroyArgs = []string{"-c"}
)

// NewDestroyCommand returns a new instance of the destroy command for destroying and deleting a local development cluster.
func NewDestroyCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a local amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroy(c)
		},
	}
}

func destroy(c cli.Interface) error {
	return updateCluster(c, destroyArgs)
}
