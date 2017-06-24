package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
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
	flags.IntVarP(&opts.managers, "managers", "m", 3, "Initial number of manager nodes")
	flags.StringVar(&opts.name, "name", "", "Cluster Label")
	flags.BoolVarP(&opts.notifications, "notifications", "n", false, "Enable/disable server notifications (default is 'false')")
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider")
	flags.StringVarP(&opts.registration, "registration", "r", configuration.RegistrationNone, "Specify the registration policy (possible values are 'none' or 'email')")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Specify tag for cluster images (use 'local' for development)")
	flags.IntVarP(&opts.workers, "workers", "w", 2, "Initial number of worker nodes")

	// aws options
	flags.String("aws-onfailure", "", "'DO_NOTHING', 'ROLLBACK' (default), or 'DELETE")
	flags.String("aws-parameter", "", "A key-value pair to supply to the AWS template")
	flags.String("aws-region", "", "The region to use when launching the instance")
	flags.String("aws-stackname", "", "The name of the AWS stack that will be created")
	flags.Bool("aws-sync", false, "If true, block until the command finishes (default: false)")
	flags.String("aws-template", "", "UNSUPPORTED: cloud formation template url")

	return cmd
}

// Map cli cluster flags to target bootstrap cluster command flags and update the cluster
func create(c cli.Interface, cmd *cobra.Command) error {
	var args []string
	var env map[string]string

	// until we finish refactoring, the local provider uses a deploy script that requires
	// remapping cli options to script args and env vars
	if opts.provider == "local" {
		// This is a map from cli cluster flag name to bootstrap script flag name
		m := map[string]string{
			"workers":       "-w",
			"managers":      "-m",
			"provider":      "-t",
			"name":          "-l",
			"tag":           "-T",
			"registration":  "-r",
			"notifications": "-n",
		}

		// TODO: only supporting local cluster management for this release
		// the following ensures that flags are added before the final command arg
		// TODO: refactor reflag to handle this
		args = []string{"bin/deploy"}
		args = reflag(cmd, m, args)
		env = map[string]string{"TAG": opts.tag, "REGISTRATION": opts.registration, "NOTIFICATIONS": strconv.FormatBool(opts.notifications)}
	} else {
		args = append(args, "init")
	}

	// plugin options should be prefixed by the name of the provider
	// for example: provider "aws" => --aws-region "us-west-2"
	cmd.Flags().Visit(func (f *flag.Flag) {
		if strings.HasPrefix(f.Name, opts.provider) {
			 // TODO: wip, don't think we will need to keep this next line
			opts.options[f.Name] = f.Value.String()

			name := strings.TrimPrefix(f.Name, opts.provider + "-")
			opt := fmt.Sprintf("--%s=%s", name, f.Value.String())
			args = append(args, opt)
			fmt.Println(opt)
		}
	})

	config := PluginConfig{
		Provider: opts.provider,
		Options: opts.options,
		DockerOpts: opts.docker,
	}
	p, err := NewPlugin(config)
	if err != nil {
		return err
	}

	return p.Run(c, args, env)
}

