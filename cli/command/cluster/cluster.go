package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/exec"
	"github.com/spf13/cobra"
)

type clusterOpts struct {
	managers int
	workers  int
	provider string
	name     string
}

var (
	opts    = &clusterOpts{3, 2, "local", ""}
	flagMap map[string]string
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
	cmd.AddCommand(NewDestroyCommand(c))
	cmd.AddCommand(NewUpdateCommand(c))
	cmd.AddCommand(NewStatusCommand(c))
	return cmd
}

func updateCluster(c cli.Interface, args []string) error {
	return exec.Run(c, "bootstrap", args)
}

// Map cli cluster flags to target bootstrap cluster command flags,
// append to and return args array
func reflag(cmd *cobra.Command, flags map[string]string, args []string) []string {
	// transform src flags to target flags and add flag and value to cargs
	for s, t := range flags {
		if cmd.Flag(s).Changed {
			args = append(args, t, cmd.Flag(s).Value.String())
		}
	}
	return args
}
