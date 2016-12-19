package main

import (
	"fmt"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
	"time"
)

// NodeLsCmd is the node ls command
var NodeLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list cluster nodes",
	Long:  `list cluster nodes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.nodeMonitor(cmd, args)
	},
}

func init() {
	NodeCmd.AddCommand(NodeLsCmd)
	NodeLsCmd.Flags().BoolP("follow", "f", false, "follow node list")
	NodeLsCmd.Flags().Bool("more", false, "more information")
}

func (g *clusterClient) nodeMonitor(cmd *cobra.Command, args []string) error {
	more := false
	if cmd.Flag("more").Value.String() == "true" {
		more = true
	}
	follow := false
	if cmd.Flag("follow").Value.String() == "true" {
		follow = true
	}
	req := &servergrpc.GetNodesInfoRequest{
		More: more,
	}

	nodeList, err := g.client.GetNodesInfo(g.ctx, req)
	if err != nil {
		return err
	}
	g.followClearScreen(follow)
	for {
		g.followMoveCursorHome(follow)
		title := "ID\tROLE\tHOSTNAME\tADDRESS\tVERSION\tSTATUS\tCPU\tMEM"
		if more {
			title += "\tAGENTID\tOS\tARCHI"
		}
		lines := []string{}
		for _, node := range nodeList.Nodes {
			if node.Error == "" {
				line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%.1f\t", node.Id, node.Role, node.Hostname, node.Address, node.DockerVersion, node.Status, node.Cpu, float64(node.Memory)/1024000000)
				if more {
					line += fmt.Sprintf("%s\t%s\t%s", node.AgentId, node.HostOs, node.HostArchitecture)
				}
				lines = append(lines, line)
			} else {
				lines = append(lines, fmt.Sprintf("Node %s Error: %s", node.Id, node.Error))
			}
		}
		g.displayInOrder("", title, lines)
		if !follow {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}
