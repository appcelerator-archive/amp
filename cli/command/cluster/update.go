package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewUpdateCommand returns a new instance of the update command for updating local development cluster.
func NewUpdateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update a local amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return update(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.workers, "workers", "w", 2, "Initial number of worker nodes")
	flags.IntVarP(&opts.managers, "managers", "m", 3, "Intial number of manager nodes")
	return cmd
}

func update(c cli.Interface, cmd *cobra.Command) error {
	// This is a map from cli cluster flag name to bootstrap script flag name
	m := map[string]string{
		"workers":  "-w",
		"managers": "-m",
	}

	args := []string{}
	return updateCluster(c, reflag(cmd, m, args))
}
