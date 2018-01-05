package cluster

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	grpcStatus "google.golang.org/grpc/status"
)

type nodeCleanupOptions struct {
	force bool
}

var (
	nodeCleanupOpts = &nodeCleanupOptions{}
)

// NewNodeCleanupCommand returns a new instance of the cleanup command for amp clusters.
func NewNodeCleanupCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cleanup",
		Aliases: []string{"clean"},
		Short:   "Remove amp cluster nodes in the *Down* state",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nodeCleanup(c)
		},
	}
	cmd.Flags().BoolVarP(&nodeCleanupOpts.force, "force", "f", false, "Force the node removal")
	return cmd
}

func nodeCleanup(c cli.Interface) error {
	req := &cluster.NodeCleanupRequest{
		Force: nodeCleanupOpts.force,
	}
	client := cluster.NewClusterClient(c.ClientConn())
	reply, err := client.ClusterNodeCleanup(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if len(reply.Nodes) == 0 {
		fmt.Println("Nothing to clean")
		return nil
	}
	fmt.Println("The following nodes were down and have been removed")
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tHOSTNAME\tROLE")
	for _, node := range reply.Nodes {
		fmt.Fprintf(w, "%s\t%s\t%s\n", node.Id, node.Hostname, strings.Title(node.Role))
	}
	w.Flush()
	return nil
}
