package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewStatusCommand returns a new instance of the status command for querying the state of amp cluster.
func NewStatusCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Retrieve details about an amp cluster",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(c)
		},
	}
}

func status(c cli.Interface) error {
	// TODO call api to get status
	return nil
}
