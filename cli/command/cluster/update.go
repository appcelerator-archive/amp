package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewUpdateCommand returns a new instance of the update command for updating local development cluster.
func NewUpdateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a local amp cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return update(c, cmd, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.workers, "workers", "w", 2, "Initial number of worker nodes")
	flags.IntVarP(&opts.managers, "managers", "m", 3, "Intial number of manager nodes")
	return cmd
}

func update(c cli.Interface, cmd *cobra.Command, args []string) error {
	flagMap = make(map[string]string)
	if cmd.Flag("workers").Changed {
		flagMap["-w"] = cmd.Flag("workers").Value.String()
	}
	if cmd.Flag("managers").Changed {
		flagMap["-m"] = cmd.Flag("managers").Value.String()
	}
	for k, v := range flagMap {
		args = append(args, k, v)
	}
	return updateCluster(c, args)
}
