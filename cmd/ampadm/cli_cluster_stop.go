package main

import (
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
)

// ClusterStopCmd is the cluster stop command
var ClusterStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop locally cluster services",
	Long:  `stop locally cluster services (server and agents)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.stopClusterServices(cmd, args)
	},
}

func init() {
	ClusterCmd.AddCommand(ClusterStopCmd)
}

func (g *clusterClient) stopClusterServices(cmd *cobra.Command, args []string) error {
	if !g.isLocalhostServer() {
		g.fatalc("Stop cluster services works only locally, --server option is invalid\n")
	}
	if cmd.Flag("silence").Value.String() == "true" {
		g.clusterLoader.silence = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		g.clusterLoader.verbose = true
	}
	if err := g.clusterLoader.init(g, ""); err != nil {
		g.fatalc("%v\n", err)
	}

	if g.clusterLoader.isServiceRunning("amp-cluster") {
		if err := client.initConnection(); err != nil {
			g.fatalc("%v\n", err)
		}
		req := &servergrpc.AmpRequest{ClientId: g.id, Silence: true}
		ret, err := g.client.GetAmpStatus(g.ctx, req)
		if err != nil {
			g.fatalc("%v\n", err)
		}
		if ret.Status == "running" || ret.Status == "partially running" {
			g.fatalc("AMP platform should be stopped first\n")
		}
	}
	if err := g.clusterLoader.stopClusterServices(); err != nil {
		g.fatalc("%v\n", err)
	}
	return nil
}
