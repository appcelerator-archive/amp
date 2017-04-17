package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewListCommand returns a new instance of the list command for amp clusters.
func NewListCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List deployed amp clusters",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c)
		},
	}
}

func list(c cli.Interface) error {
	// TODO: only supporting local cluster management for this release
	// TODO call api to list clusters
	c.Console().Println(DefaultLocalClusterID)
	return nil
}
