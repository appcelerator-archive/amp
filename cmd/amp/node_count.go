package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
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
		return containerCount(AMP, cmd, args)
	},
}

func init() {
	NodeCmd.AddCommand(NodeCountCmd)
	NodeCountCmd.Flags().BoolP("follow", "f", false, "follow node list")
	NodeCountCmd.Flags().StringP("node", "n", "", "specify the node onto apply the purge, default all")
}

func containerCount(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if err := amp.ConnectAdmServer(); err != nil {
		manager.fatalf("ampadmin_server is not available, start stack ampadmin first: 'amp pf start ampadmin':\n")
	}
	client := servergrpc.NewClusterServerServiceClient(amp.ConnAdmServer)
	req := &servergrpc.GetNodesInfoRequest{}
	follow := false
	req.Node = cmd.Flag("node").Value.String()
	if cmd.Flag("follow").Value.String() == "true" {
		follow = true
	}
	manager.followClearScreen(follow)
	for {
		nodeList, err := client.GetNodesInfo(ctx, req)
		if err != nil {
			manager.fatalf("%v\n", err)
		}
		manager.followMoveCursorHome(follow)
		list := []string{}
		for _, node := range nodeList.Nodes {
			list = append(list, fmt.Sprintf("%s\t%d\t%d\t%d\t%d\t%d", node.NodeName, node.NbContainers, node.NbContainersRunning, node.NbContainersPaused, node.NbContainersStopped, node.Images))
		}
		manager.displayInOrder("NODE\tCONTAINER\tCONTAINER\tCONTAINER\tCONTAINER\tIMAGE", "HOSTNAME (ID)\tTOTAL\tRUNNING\tPAUSED\tSTOPPED\tTOTAL", list)
		if !follow {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}
