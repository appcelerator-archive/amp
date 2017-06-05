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
	flags.StringVarP(&opts.tag, "tag", "t", "0.12.0", "Specify tag for bootstrap images (default is '0.12.0', use 'local' for development)")
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVar(&opts.name, "name", "", "Cluster Label")
	return cmd
}

func remove(c cli.Interface, cmd *cobra.Command) error {
	// This is a map from cli cluster flag name to bootstrap script flag name
	m := map[string]string{
		"provider": "-t",
		"tag":      "-T",
	}
	// TODO: only supporting local cluster management for this release
	args := []string{"bin/deploy", "-d"}
	args = reflag(cmd, m, args)
	env := map[string]string{"TAG": opts.tag}
	return queryCluster(c, args, env)
}
