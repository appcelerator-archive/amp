package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

var (
	createArgs = []string{"-p", "docker"}
)

// NewCreateCommand returns a new instance of the create command for bootstrapping a local development cluster.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a local amp cluster",
		Long:  `The create command initializes a local amp cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, args)
		},
	}
}

func create(c cli.Interface, args []string) error {
	return updateCluster(c, append(createArgs[:], args[:]...))
}
