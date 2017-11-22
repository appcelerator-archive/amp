package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewRemoveCommand returns a new instance of the remove command for destroying a cluster.
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Aliases: []string{"remove", "destroy"},
		Short:   "Destroy an amp cluster",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Specify tag for cluster plugin image")
	flags.StringVar(&opts.name, "name", "", "Cluster Label")

	// local options
	flags.Bool("local-force-leave", false, "Force leave the swarm")

	// aws options
	flags.String("aws-region", "", "Region to use to delete the instance")
	flags.String("aws-stackname", "", "Name of the AWS stack to be deleted")
	flags.Bool("aws-sync", false, "If true, block until the command finishes (default: false)")
	flags.String("aws-access-key-id", "", "aws credential: access key id")
	flags.String("aws-secret-access-key", "", "aws credential: secret access key")

	return cmd
}

func remove(c cli.Interface, cmd *cobra.Command) error {
	return runPluginCommand(c, cmd, "destroy")
}
