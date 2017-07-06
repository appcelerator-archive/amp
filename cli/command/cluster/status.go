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
	//flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Specify tag for the amp-bootstrap image")
	return cmd
}

func status(c cli.Interface, cmd *cobra.Command) error {
	// This is a map from cli cluster flag name to bootstrap script flag name
	m := map[string]string{
		//"provider": "-t",
		"tag": "-T",
	}
	// TODO call api to get status
	args := []string{"bin/deploy", "-s"}
	args = reflag(cmd, m, args)
	env := map[string]string{"TAG": opts.tag}
	status := queryCluster(c, args, env)
	if status != nil {
		c.Console().Println("cluster status: not running")
	} else {
		c.Console().Println("cluster: running")
	}
	return nil
}
