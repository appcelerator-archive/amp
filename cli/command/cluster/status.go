package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewStatusCommand returns a new instance of the status command for by providing groups and instances of local cluster.
func NewStatusCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Short:   "Retrieve details about a local amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(c)
		},
	}
}

func status(c cli.Interface) error {
	// TODO call api to get status
	return nil
}
