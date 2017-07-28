package cluster

import (
	"fmt"
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
	// TODO: not clear yet if we'll need this
	options map[string]string
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
		tag:           "latest",
		registration:  configuration.RegistrationDefault,
		notifications: true,
		// TODO: not clear yet if we'll need this
		options: map[string]string{},
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
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	cmd.AddCommand(NewStatusCommand(c))
	cmd.AddCommand(NewUpdateCommand(c))
	cmd.AddCommand(NewNodeCommand(c))
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
