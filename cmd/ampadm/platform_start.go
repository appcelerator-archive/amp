package main

import (
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
)

// PlatformStart is the platform start command
var PlatformStart = &cobra.Command{
	Use:   "start",
	Short: "Start platform",
	Long:  `Start all AMP platform services.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.startAMP(cmd, args)
	},
}

func init() {
	PlatformStart.Flags().BoolP("force", "f", false, "Start all possible services, do not stop on error")
	PlatformStart.Flags().BoolP("local", "l", false, "use local amp image")
	PlatformStart.Flags().String("tag", "", "use alternate tag for amp image")
	PlatformCmd.AddCommand(PlatformStart)
}

func (g *clusterClient) startAMP(cmd *cobra.Command, args []string) error {
	req := &servergrpc.AmpRequest{}
	if cmd.Flag("silence").Value.String() == "true" {
		req.Silence = true
		g.clusterLoader.silence = true
	}
	if cmd.Flag("force").Value.String() == "true" {
		req.Force = true
		g.clusterLoader.force = true
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

	if err := g.initConnection(); err != nil {
		g.fatalc("%v\n", err)
	}
	req.ClientId = g.id
	if _, err := g.client.AmpStart(g.ctx, req); err != nil {
		g.fatalc("%v\n", err)
	}
	return nil
}
