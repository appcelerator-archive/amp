package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewStatusCommand returns a new instance of the status command for querying the state of amp cluster.
func NewStatusCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Retrieve details about an amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Specify tag for cluster plugin image")
	return cmd
}

func status(c cli.Interface, cmd *cobra.Command) error {
	return runPluginCommand(c, cmd, "info")
}
