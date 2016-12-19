package main

import (
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
	"time"
)

// PlatformMonitor is the platform monitorcommand
var PlatformMonitor = &cobra.Command{
	Use:   "monitor",
	Short: "Display AMP platform services",
	Long:  `Display AMP platform services information and states`,
	Run: func(cmd *cobra.Command, args []string) {
		client.ampMonitor(cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformMonitor)
	PlatformMonitor.Flags().BoolP("local", "l", false, "use local amp image to start locally cluster server/agent if not started")
}

func (g *clusterClient) ampMonitor(cmd *cobra.Command, args []string) error {
	req := &servergrpc.AmpRequest{Silence: true}
	if cmd.Flag("local").Value.String() == "true" {
		g.clusterLoader.local = true
	}
	//---Treatement only for local usage
	if g.isLocalhostServer() {
		g.startLocalClusterService()
	}
	//---
	if err := client.initConnection(); err != nil {
		g.fatalc("%v\n", err)
	}
	req.ClientId = g.id
	g.verbose = true
	g.followClearScreen(true)
	for {
		lineList, err := g.client.AmpMonitor(g.ctx, req)
		if err != nil {
			g.fatalc("%v\n", err)
		}
		g.followMoveCursorHome(true)
		for _, output := range lineList.Outputs {
			g.printfc(int(output.OutputType), "%s\n", output.Output)
		}
		time.Sleep(1 * time.Second)
	}
}
