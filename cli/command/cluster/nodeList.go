package cluster

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type nodeListOptions struct {
	quiet bool
}

var (
	nodeListopts = &nodeListOptions{}
)

// NewNodeListCommand returns a new instance of the list command for amp clusters.
func NewNodeListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List deployed amp cluster nodes",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nodeList(c)
		},
	}
	cmd.Flags().BoolVarP(&nodeListopts.quiet, "quiet", "q", false, "Only display the node id")
	return cmd
}

func nodeList(c cli.Interface) error {
	req := &cluster.NodeListRequest{}
	client := cluster.NewClusterClient(c.ClientConn())
	reply, err := client.NodeList(context.Background(), req)
	if err != nil {
		return errors.New(grpc.ErrorDesc(err))
	}
	if !nodeListopts.quiet {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tHOSTNAME\tSTATUS\tAVAILABILITY\tMANAGER STATUS")
		for _, node := range reply.Nodes {
			leader := ""
			if node.ManagerLeader {
				leader = "Leader"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", node.Id, node.Hostname, node.Status, node.Availability, leader)
		}
		w.Flush()
	} else {
		for _, node := range reply.Nodes {
			c.Console().Println(node.Id)
		}
	}
	return nil
}
