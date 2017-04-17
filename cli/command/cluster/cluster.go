package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

type clusterOpts struct {
	managers int
	workers  int
	provider string
	name     string
}

const (
	DefaultLocalClusterID = "f573e897-7aa0-4516-a195-42ee91039e97"
)

var (
	opts = &clusterOpts{3, 2, "local", ""}
)

// NewClusterCommand returns a new instance of the cluster command.
func NewClusterCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Cluster management operations",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			c.Console().Infoln("Note: only 'local' cluster provider supported in this release")
		},
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	cmd.AddCommand(NewUpdateCommand(c))
	cmd.AddCommand(NewStatusCommand(c))
	cmd.AddCommand(NewListCommand(c))
	return cmd
}

func queryCluster(c cli.Interface, args []string) error {
	err := Run(c, args)
	if err != nil {
		// TODO: the local cluster is the only one that can be managed this release
		c.Console().Println(DefaultLocalClusterID)
	}
	return err
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
