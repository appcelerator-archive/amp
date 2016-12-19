package main

import (
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
)

// PlatformPull is the platform pull command
var PlatformPull = &cobra.Command{
	Use:   "pull",
	Short: "Pull platform images",
	Long:  `Pull all AMP platform images`,
	Run: func(cmd *cobra.Command, args []string) {
		client.platformPull(cmd, args)
	},
}

func init() {
	PlatformPull.Flags().StringP("node", "n", "", "specify the node onto apply the purge, default all")
	PlatformPull.Flags().BoolP("local", "l", false, "use local amp image to start locally cluster server/agent if not started")
	PlatformPull.Flags().String("tag", "", "use alternate tag for amp image")
	PlatformCmd.AddCommand(PlatformPull)
}

func (g *clusterClient) platformPull(cmd *cobra.Command, args []string) error {
	req := &servergrpc.AmpRequest{}
	req.Node = cmd.Flag("node").Value.String()
	if cmd.Flag("silence").Value.String() == "true" {
		req.Silence = true
		g.clusterLoader.silence = true
		g.clusterLoader.verbose = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		req.Verbose = true
		g.clusterLoader.verbose = true
	}
	if cmd.Flag("local").Value.String() == "true" {
		req.Local = true
		g.clusterLoader.local = true

	}

	g.clusterLoader.ampTag = cmd.Flag("tag").Value.String()
	//---Treatement only for local usage
	if g.isLocalhostServer() {
		g.startLocalClusterService()
	}
	//---

	if err := client.initConnection(); err != nil {
		g.fatalc("%v\n", err)
	}
	req.ClientId = g.id

	if _, err := g.client.AmpPull(g.ctx, req); err != nil {
		g.fatalc("%v\n", err)
	}
	return nil
}
