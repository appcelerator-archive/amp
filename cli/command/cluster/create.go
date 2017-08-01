package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
)

// NewCreateCommand returns a new instance of the create command for bootstrapping an cluster.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Set up a cluster in swarm mode",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.name, "name", "", "Cluster Label")
	flags.BoolVarP(&opts.notifications, "notifications", "n", false, "Enable/disable server notifications (default is 'false')")
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider (\"local\" (default) or \"aws\")")
	flags.StringVarP(&opts.registration, "registration", "r", configuration.RegistrationNone, "Specify the registration policy (possible values are 'none' or 'email')")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Specify tag for cluster plugin image")

	// local options
	flags.String("local-advertise-addr", "", "Swarm advertised address (format: <ip|interface>[:port])")
	flags.Bool("local-force-new-cluster", false, "Force initialization of a new swarm")
	flags.Bool("local-fast", false, "Skip cluster service smoke tests")
	flags.Bool("local-no-monitoring", false, "Don't deploy monitoring services")

	// aws options
	flags.String("aws-onfailure", "", "'DO_NOTHING', 'ROLLBACK' (default), or 'DELETE")
	flags.StringSlice("aws-parameter", []string{}, "Key-value pairs to pass through to the AWS CloudFormation template")
	flags.String("aws-region", "", "The region to use when launching the instance")
	flags.String("aws-stackname", "", "The name of the AWS stack that will be created")
	flags.Bool("aws-sync", false, "If true, block until the command finishes (default: false)")
	flags.String("aws-template", "", "UNSUPPORTED: cloud formation template url")

	return cmd
}

func create(c cli.Interface, cmd *cobra.Command) error {
	return runPluginCommand(c, cmd, "init")
}
