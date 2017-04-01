package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewCreateCommand returns a new instance of the create command for bootstrapping a local development cluster.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a local amp cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.workers, "workers", "w", 2, "Initial number of worker nodes")
	flags.IntVarP(&opts.managers, "managers", "m", 3, "Intial number of manager nodes")
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVar(&opts.name, "name", "", "Cluster Label")
	return cmd
}

func create(c cli.Interface, cmd *cobra.Command, args []string) error {
	flagMap = make(map[string]string)
	if cmd.Flag("workers").Changed {
		flagMap["-w"] = cmd.Flag("workers").Value.String()
	}
	if cmd.Flag("managers").Changed {
		flagMap["-m"] = cmd.Flag("managers").Value.String()
	}
	if cmd.Flag("provider").Changed {
		flagMap["-t"] = cmd.Flag("provider").Value.String()
	}
	if cmd.Flag("name").Changed {
		flagMap["-l"] = cmd.Flag("name").Value.String()
	}
	for k, v := range flagMap {
		args = append(args, k, v)
	}
	return updateCluster(c, args)
}
