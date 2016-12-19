package main

import (
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
)

// PlatformStatus is the platform status command
var PlatformStatus = &cobra.Command{
	Use:   "status",
	Short: "Get AMP platform status",
	Long:  `Get AMP platform global status (stopped, partially running, running command return 1 if status is not running`,
	Run: func(cmd *cobra.Command, args []string) {
		client.getAMPStatus(cmd, args)
	},
}

func init() {
	PlatformStatus.Flags().BoolP("silence", "s", false, "no console output at all")
	PlatformStatus.Flags().BoolP("local", "l", false, "use local amp image")
	PlatformCmd.AddCommand(PlatformStatus)
}

func (g *clusterClient) getAMPStatus(cmd *cobra.Command, args []string) error {

	req := &servergrpc.AmpRequest{}
	if cmd.Flag("silence").Value.String() == "true" {
		req.Silence = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		req.Verbose = true
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

	_, err := g.client.GetAmpStatus(g.ctx, req)
	if err != nil {
		g.fatalc("%v\n", err)
	}
	return nil
}
