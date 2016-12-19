package main

import (
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
)

// PlatformStop is the platform stop command
var PlatformStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop platform",
	Long:  `Stop all AMP platform services.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.stopAMP(cmd, args)
	},
}

func init() {
	PlatformStop.Flags().BoolP("silence", "s", false, "no console output at all")
	PlatformCmd.AddCommand(PlatformStop)
}

func (g *clusterClient) stopAMP(cmd *cobra.Command, args []string) error {
	req := &servergrpc.AmpRequest{}
	if err := client.initConnection(); err != nil {
		g.fatalc("%v\n", err)
	}
	if cmd.Flag("silence").Value.String() == "true" {
		req.Silence = true
		g.clusterLoader.silence = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		req.Verbose = true
		g.clusterLoader.verbose = true
	}

	if err := g.initConnection(); err != nil {
		g.fatalc("%v\n", err)
	}
	req.ClientId = g.id
	if _, err := g.client.AmpStop(g.ctx, req); err != nil {
		g.fatalc("%v\n", err)
	}
	//---Treatement only for local usage
	if g.isLocalhostServer() {
		if err := g.stopClusterServices(cmd, args); err != nil {
			g.fatalc("%v\n", err)
		}
	}
	//---
	return nil
}
