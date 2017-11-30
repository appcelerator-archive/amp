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

type docker struct {
	Volumes []string
}

type clusterOpts struct {
	docker
	managers      int
	workers       int
	provider      string
	name          string
	tag           string
	registration  string
	notifications bool
}

var (
	opts = &clusterOpts{
		docker: docker{
			Volumes: []string{},
		},
		managers:      3,
		workers:       2,
		provider:      "local",
		name:          "",
		tag:           "0.17.0",
		registration:  configuration.RegistrationDefault,
		notifications: true,
	}
)

// NewClusterCommand returns a new instance of the cluster command.
func NewClusterCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Cluster management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}

	cmd.PersistentFlags().StringSliceVarP(&opts.Volumes, "volume", "v", []string{}, "Bind mount a volume")

	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	cmd.AddCommand(NewStatusCommand(c))
	cmd.AddCommand(NewNodeCommand(c))
	// Update command is not fully implemented, disabled for now
	//cmd.AddCommand(NewUpdateCommand(c))
	return cmd
}

// Strip provider prefix from CLI options to provider-specific options that will be supplied to the plugin.
// For example, `--aws-region us-west-2` becomes `--region us-west-2` (as expected by the aws plugin)
func stripPrefixes(cmd *cobra.Command, args []string) []string {
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
	return args
}

func runPluginCommand(c cli.Interface, cmd *cobra.Command, command string) error {
	// args and env will be supplied to the cluster plugin container
	var args []string
	env := map[string]string{}
	if opts.tag != "" {
		env["TAG"] = opts.tag
	}
	if opts.registration != "" {
		env["REGISTRATION"] = opts.registration
	}
	if opts.notifications {
		env["NOTIFICATIONS"] = strconv.FormatBool(opts.notifications)
	}

	args = append(args, command)
	args = stripPrefixes(cmd, args)

	config := PluginConfig{
		Provider:   opts.provider,
		DockerOpts: opts.docker,
	}

	p, err := NewPlugin(config)
	if err != nil {
		return err
	}

	return p.Run(c, args, env)
}
