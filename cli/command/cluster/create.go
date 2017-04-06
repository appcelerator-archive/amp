package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewCreateCommand returns a new instance of the create command for bootstrapping a local development cluster.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a local amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.workers, "workers", "w", 2, "Initial number of worker nodes")
	flags.IntVarP(&opts.managers, "managers", "m", 3, "Intial number of manager nodes")
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVar(&opts.name, "name", "", "Cluster Label")
	return cmd
}

// Map cli cluster flags to target bootstrap cluster command flags and update the cluster
func create(c cli.Interface, cmd *cobra.Command) error {
	// This is a map from cli cluster flag name to bootstrap script flag name
	m := map[string]string{
		"workers":  "-w",
		"managers": "-m",
		"provider": "-t",
		"name":     "-l",
	}

	args := []string{}
	return updateCluster(c, reflag(cmd, m, args))
}
