package main

import (
	"fmt"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
	"time"
)

// NodeCountCmd is the node couunt command
var NodeCountCmd = &cobra.Command{
	Use:   "count",
	Short: "count containers/images in nodes",
	Long:  `count containers/images in nodes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.containerCount(cmd, args)
	},
}

func init() {
	NodeCmd.AddCommand(NodeCountCmd)
	NodeCountCmd.Flags().BoolP("follow", "f", false, "follow node list")
	NodeCountCmd.Flags().StringP("node", "n", "", "specify the node onto apply the purge, default all")
}

func (g *clusterClient) containerCount(cmd *cobra.Command, args []string) error {
	req := &servergrpc.GetNodesInfoRequest{}
	follow := false
	req.Node = cmd.Flag("node").Value.String()
	if cmd.Flag("follow").Value.String() == "true" {
		follow = true
	}
	g.followClearScreen(follow)
	for {
		nodeList, err := g.client.GetNodesInfo(g.ctx, req)
		if err != nil {
			g.fatalc("%v\n", err)
		}
		g.followMoveCursorHome(follow)
		list := []string{}
		for _, node := range nodeList.Nodes {
			list = append(list, fmt.Sprintf("%s\t%d\t%d\t%d\t%d\t%d", node.NodeName, node.NbContainers, node.NbContainersRunning, node.NbContainersPaused, node.NbContainersStopped, node.Images))
		}
		g.displayInOrder("NODE\tCONTAINER\tCONTAINER\tCONTAINER\tCONTAINER\tIMAGE", "HOSTNAME (ID)\tTOTAL\tRUNNING\tPAUSED\tSTOPPED\tTOTAL", list)
		if !follow {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}
