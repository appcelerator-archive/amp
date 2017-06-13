package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
)

type clusterOpts struct {
	log           int
	name          string
	tag           string
	provider      string
	managers      int
	workers       int
	registration  string
	notifications bool
	region        string
	domain        string
}

var (
	opts = &clusterOpts{
		log:           4,
		name:          "",
		tag:           "latest",
		provider:      "local",
		managers:      3,
		workers:       2,
		registration:  configuration.RegistrationDefault,
		notifications: true,
		region:        "",
		domain:        "",
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
	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	cmd.AddCommand(NewStatusCommand(c))
	cmd.AddCommand(NewUpdateCommand(c))
	cmd.AddCommand(NewNodeCommand(c))
	return cmd
}

func queryCluster(c cli.Interface, args []string, env map[string]string) error {
	if err := check(opts.provider); err != nil {
		return err
	}
	err := Run(c, args, env)
	return err
}

// Map cli cluster flags to target bootstrap cluster command flags,
// append to and return args array
func reflag(cmd *cobra.Command, flags map[string]string, args []string) []string {
	// transform src flags to target flags and add flag and value to cargs
	for s, t := range flags {
		// Allow default values to be passed from cli
		if cmd.Flag(s).Value.String() != "" {
			args = append(args, t, cmd.Flag(s).Value.String())
		}
	}
	return args
}
