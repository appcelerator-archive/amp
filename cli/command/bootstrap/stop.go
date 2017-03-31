package bootstrap

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

var (
	stopArgs = []string{"-c"}
)

// NewStopCommand returns a new instance of the stop command for stopping and deleting a local development cluster.
func NewStopCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop a local amp cluster",
		Long:  `The stop command stops and cleans up a local amp cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stop(c, args)
		},
	}
}

func stop(c cli.Interface, args []string) error {
	return updateCluster(append(stopArgs[:], args[:]...))
}
