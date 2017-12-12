package cluster

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/cluster"
	"github.com/appcelerator/amp/cli"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
	grpcStatus "google.golang.org/grpc/status"
)

type nodeListOptions struct {
	quiet       bool
	nodeLabel   string
	engineLabel string
}

var (
	nodeListopts = &nodeListOptions{}
)

// NewNodeListCommand returns a new instance of the list command for amp clusters.
func NewNodeListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List amp cluster nodes",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nodeList(c)
		},
	}
	cmd.Flags().BoolVarP(&nodeListopts.quiet, "quiet", "q", false, "Only display the node id")
	cmd.Flags().StringVar(&nodeListopts.nodeLabel, "label", "", "Filter nodes with a node label")
	cmd.Flags().StringVar(&nodeListopts.engineLabel, "engine-label", "", "Filter nodes with an engine label")
	return cmd
}

// ByRole implements sort.Interface for []NodeReply based on the Role field.
type ByRole []*cluster.NodeReply

func (a ByRole) Len() int           { return len(a) }
func (a ByRole) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRole) Less(i, j int) bool { return a[i].Role < a[j].Role }

func nodeList(c cli.Interface) error {
	req := &cluster.NodeListRequest{
		NodeLabel:   nodeListopts.nodeLabel,
		EngineLabel: nodeListopts.engineLabel,
	}
	client := cluster.NewClusterClient(c.ClientConn())
	reply, err := client.ClusterNodeList(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if nodeListopts.quiet {
		for _, node := range reply.Nodes {
			c.Console().Println(node.Id)
		}
		return nil
	}
	if len(reply.Nodes) == 0 {
		fmt.Println("No node found")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	if len(reply.Nodes) == 1 {
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
	} else {
		fmt.Fprintln(w, "ID\tHOSTNAME\tSTATUS\tAVAILABILITY\tROLE\tVERSION\tCPUS\tMEMORY")
		sort.Sort(ByRole(reply.Nodes))
		for _, node := range reply.Nodes {
			role := strings.Title(node.Role)
			if node.Role == "manager" && node.ManagerLeader {
				role = "Manager (leader)"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%0.2f\t%s\n", node.Id, node.Hostname, strings.Title(node.Status), strings.Title(node.Availability), role, node.EngineVersion, float64(node.NanoCpus)/1000000000, units.BytesSize(float64(node.MemoryBytes)))
		}
	}
	w.Flush()
	return nil
}
