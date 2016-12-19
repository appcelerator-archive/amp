package main

import (
	"github.com/spf13/cobra"
)

// ClusterStartCmd is the cluster start command
var ClusterStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start locally cluster services",
	Long:  `start locally cluster services (server and agents)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.startClusterServices(cmd, args)
	},
}

func init() {
	ClusterCmd.AddCommand(ClusterStartCmd)
	ClusterStartCmd.Flags().BoolP("force", "f", false, "Start all possible services, do not stop on error")
	ClusterStartCmd.Flags().BoolP("local", "l", false, "use local tag for amp image")
	ClusterStartCmd.Flags().String("tag", "", "use alternate tag for amp image")
}

func (g *clusterClient) startClusterServices(cmd *cobra.Command, args []string) error {
	if !g.isLocalhostServer() {
		g.fatalc("Start cluster services works only locally, --server option is invalid\n")
	}
	if cmd.Flag("silence").Value.String() == "true" {
		g.clusterLoader.silence = true
	}
	if cmd.Flag("force").Value.String() == "true" {
		g.clusterLoader.force = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		g.clusterLoader.verbose = true
	}
	if cmd.Flag("local").Value.String() == "true" {
		g.clusterLoader.local = true
	}
	g.clusterLoader.ampTag = cmd.Flag("tag").Value.String()
	if err := g.clusterLoader.init(g, ""); err != nil {
		g.fatalc("%v\n", err)
	}
	if err := g.clusterLoader.startClusterServices(); err != nil {
		g.fatalc("%v\n", err)
	}
	return nil
}
