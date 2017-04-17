package cluster

import (
	"errors"
	"strings"

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
		Use:   "cluster",
		Short: "Cluster management operations",
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// TODO special case handling for cluster this release
			local := strings.HasPrefix(c.Server(), "127.0.0.1") ||
				strings.HasPrefix(c.Server(), "localhost")
			if !local {
				return errors.New("Note: only cluster operations with '--server=localhost' supported in this release")
			}
			return nil
		},
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	cmd.AddCommand(NewStatusCommand(c))
	cmd.AddCommand(NewUpdateCommand(c))
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
