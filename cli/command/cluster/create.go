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
	// NOTE: the top level options are in a state of transition right now. The focus is on the aws plugin
	// during this refactoring.
	flags.StringVar(&opts.name, "name", "", "Cluster Label")
	flags.BoolVarP(&opts.notifications, "notifications", "n", false, "Enable/disable server notifications (default is 'false')")
	flags.StringVar(&opts.provider, "provider", "local", "Cluster provider (\"local\" (default) or \"aws\")")
	flags.StringVarP(&opts.registration, "registration", "r", configuration.RegistrationNone, "Specify the registration policy (possible values are 'none' or 'email')")
	flags.StringVarP(&opts.tag, "tag", "t", c.Version(), "Specify tag for cluster images (use 'local' for development)")

	// local options
	flags.String("local-managers", "1", "Initial number of local manager nodes")
	flags.String("local-workers", "2", "Initial number of local worker nodes")

	// aws options
	flags.String("aws-onfailure", "", "'DO_NOTHING', 'ROLLBACK' (default), or 'DELETE")
	flags.StringSlice("aws-parameter", []string{}, "Key-value pairs to pass through to the AWS CloudFormation template")
	flags.String("aws-region", "", "The region to use when launching the instance")
	flags.String("aws-stackname", "", "The name of the AWS stack that will be created")
	flags.Bool("aws-sync", false, "If true, block until the command finishes (default: false)")
	flags.String("aws-template", "", "UNSUPPORTED: cloud formation template url")

	return cmd
}

// Map cli cluster flags to target bootstrap cluster command flags and update the cluster
func create(c cli.Interface, cmd *cobra.Command) error {
	// args and env will be supplied to the cluster plugin container
	var args []string
	var env map[string]string

	// NOTE: all the local implementation is in transition right now -- I'll try to keep it working
	// through the refactoring to support aws, but right now aws is the priority and I'll have to
	// circle back to local after.
	// until we finish refactoring, the local provider uses a deploy script that requires
	// remapping cli options to script args and env vars
	if opts.provider == "local" {
		// This is a map from cli cluster flag name to bootstrap script flag name
		m := map[string]string{
			"local-workers":  "-w",
			"local-managers": "-m",
			"provider":       "-t",
			"name":           "-l",
			"tag":            "-T",
			"registration":   "-r",
			"notifications":  "-n",
		}

		// the following ensures that flags are added before the final command arg
		// TODO: refactor reflag to handle this
		args = []string{"bin/deploy"}
		args = reflag(cmd, m, args)
		env = map[string]string{"TAG": opts.tag, "REGISTRATION": opts.registration, "NOTIFICATIONS": strconv.FormatBool(opts.notifications)}
	} else {
		args = append(args, "init")
	}

	// Strip provider prefix from CLI options to provider-specific options that will be supplied to the plugin.
	// For example, `--aws-region us-west-2` becomes `--region us-west-2` (as expected by the aws plugin)
	cmd.Flags().Visit(func(f *flag.Flag) {
		if strings.HasPrefix(f.Name, opts.provider) {
			name := strings.TrimPrefix(f.Name, opts.provider+"-")

			if strings.HasSuffix(f.Value.Type(), "Slice") {
				// error shouldn't happen here since we're asking for the slice it says it has
				slice, _ := cmd.Flags().GetStringSlice(f.Name)
				for _, val := range slice {
					opt := fmt.Sprintf("--%s", name)
					args = append(args, opt, val)
				}

			} else {
				opt := fmt.Sprintf("--%s", name)
				args = append(args, opt, f.Value.String())
			}

		}
	})

	config := PluginConfig{
		Provider:   opts.provider,
		DockerOpts: opts.docker,
		// TODO: not clear yet if we'll need this
		Options: opts.options,
	}

	p, err := NewPlugin(config)
	if err != nil {
		return err
	}

	return p.Run(c, args, env)
}
