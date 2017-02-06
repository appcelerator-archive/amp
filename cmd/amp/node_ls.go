package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
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
		return nodeMonitor(AMP, cmd, args)
	},
}

func init() {
	NodeCmd.AddCommand(NodeLsCmd)
	NodeLsCmd.Flags().BoolP("follow", "f", false, "follow node list")
	NodeLsCmd.Flags().Bool("more", false, "more information")
}

func nodeMonitor(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if err := amp.ConnectAdmServer(); err != nil {
		manager.printf(colError, "%v\n", err)
		manager.fatalf("ampadmin_server is not available, start stack ampadmin first: 'amp pf start ampadmin'\n")
	}
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

	client := servergrpc.NewClusterServerServiceClient(amp.ConnAdmServer)
	nodeList, err := client.GetNodesInfo(ctx, req)
	if err != nil {
		return err
	}
	manager.followClearScreen(follow)
	for {
		manager.followMoveCursorHome(follow)
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
		manager.displayInOrder("", title, lines)
		if !follow {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}
