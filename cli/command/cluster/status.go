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
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(c, cmd)
		},
	}
}

func status(c cli.Interface, cmd *cobra.Command) error {
	// TODO call api to get status
	args := []string{"bootstrap/bootstrap", "-s", DefaultLocalClusterID}
	status := queryCluster(c, args, nil)
	if status != nil {
		c.Console().Println("cluster status: not running")
	} else {
		c.Console().Println("cluster: running")
	}
	return nil
}
