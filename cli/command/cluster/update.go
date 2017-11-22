package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewUpdateCommand returns a new instance of the update command for updating an cluster.
func NewUpdateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [OPTIONS]",
		Short:   "Update an amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return update(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Specify tag for cluster plugin image")
	flags.String("aws-access-key-id", "", "aws credential: access key id")
	flags.String("aws-secret-access-key", "", "aws credential: secret access key")
	return cmd
}

func update(c cli.Interface, cmd *cobra.Command) error {
	return runPluginCommand(c, cmd, "update")
}
