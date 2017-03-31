package bootstrap

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

var (
	startArgs = []string{"-p", "docker"}
)

// NewStartCommand returns a new instance of the start command for bootstrapping a local development cluster.
func NewStartCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start a local amp cluster",
		Long:  `The start command initializes a local amp cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(args)
		},
	}
}

func start(args []string) error {
	return updateCluster(append(startArgs[:], args[:]...))
}
