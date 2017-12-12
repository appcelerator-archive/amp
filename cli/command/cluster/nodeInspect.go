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
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
	grpcStatus "google.golang.org/grpc/status"
)

// NewNodeInspectCommand returns a new instance of the node inspect command for amp clusters.
func NewNodeInspectCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect NODE_ID",
		Short:   "Inspect an amp cluster node",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nodeInspect(c, args[0])
		},
	}
	return cmd
}

func nodeInspect(c cli.Interface, id string) error {
	req := &cluster.NodeListRequest{
		Id: id,
	}
	client := cluster.NewClusterClient(c.ClientConn())
	reply, err := client.ClusterNodeList(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if len(reply.Nodes) == 0 {
		fmt.Println("No such node:", id)
		return nil
	}
	if len(reply.Nodes) == 1 {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		node := reply.Nodes[0]
		fmt.Fprintf(w, "ID\t%s\n", node.Id)
		fmt.Fprintf(w, "Name\t%s\n", node.Hostname)
		fmt.Fprintf(w, "Status\t%s\n", strings.Title(node.Status))
		fmt.Fprintf(w, "Availability\t%s\n", strings.Title(node.Availability))
		fmt.Fprintf(w, "Role\t%s", strings.Title(node.Role))
		if node.Role == "manager" && node.ManagerLeader {
			fmt.Fprintf(w, " (leader)\n")
		} else {
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "Engine Version\t%s\n", node.EngineVersion)
		fmt.Fprintf(w, "CPUs\t%0.2f\n", float64(node.NanoCpus)/1000000000)
		fmt.Fprintf(w, "Memory\t%s\n", units.BytesSize(float64(node.MemoryBytes)))
		fmt.Fprintf(w, "Node Labels\n")
		for k, v := range node.NodeLabels {
			fmt.Fprintf(w, "\t- %s=%s\n", k, v)
		}
		fmt.Fprintf(w, "Engine Labels\n")
		for k, v := range node.EngineLabels {
			fmt.Fprintf(w, "\t- %s=%s\n", k, v)
		}
		fmt.Fprintf(w, "Engine Plugins\n")
		for _, p := range node.EnginePlugins {
			fmt.Fprintf(w, "\t- Type = %s, Name = %s\n", p.Type, p.Name)
		}
		w.Flush()
	}
	return nil

}
