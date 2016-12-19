package main

import (
	"github.com/spf13/cobra"
)

var bootstrapCmd = &cobra.Command{
	Use:     "bootstrap",
	Short:   "Bootstraps the node",
	Long:    `Installs Docker, setups Swarm, run the AMP agent`,
	Aliases: []string{"bs"},
	Run: func(cmd *cobra.Command, args []string) {
		client.bootstrap(cmd, args)
	},
}

// Flag variables
var swarmManagerHost string
var swarmManagerPort string
var swarmJoinToken string
var ampImageTag string
var createSwarm bool
var manager bool

func init() {
	bootstrapCmd.Flags().BoolVarP(&manager, "manager", "m", false, "bootstraps a Swarm manager node")
	bootstrapCmd.Flags().BoolVarP(&createSwarm, "create", "c", false, "initializes the Swarm (needs option -m)")
	bootstrapCmd.Flags().StringVarP(&swarmJoinToken, "token", "", "", "Swarm join token")
	bootstrapCmd.Flags().StringVarP(&swarmManagerHost, "host", "", "", "Swarm manager hostname")
	bootstrapCmd.Flags().StringVarP(&swarmManagerPort, "port", "", "2377", "Swarm manager port")
	bootstrapCmd.Flags().StringVarP(&ampImageTag, "tag", "", Version, "AMP image version")
	RootCmd.AddCommand(bootstrapCmd)
}

// Docker daemon install/update with the official Docker script
// update the Docker configuration with:
// - ulimit for file descriptors
// - remote API
// enable and start the Docker daemon
// join a Swarm manager (or init the swarm)
// if on a manager, create:
// overlay networks amp-public and amp-infra,
// a global service with the amp cluster server
// a global service with the amp cluster agent
func (c *clusterClient) bootstrap(cmd *cobra.Command, args []string) {
	c.printfc(colInfo, "Bootstrapping the node with CLI version %s\n", Version)

	var version string

	_ = c.startDockerEngine()
	version, _ = c.getDockerServerVersion()
	ok, err := c.validateInstalledDockerEngineVersion(version)
	if err != nil {
		c.fatalc("Failed to validate Docker version: %s\n", err)
	}
	if !ok {
		if version != "0.0.0" {
			c.printfc(colInfo, "Installed version %s is lower than %s\n", version, minimumDockerEngineVersion)
		}
		if err := c.installDockerEngine(); err != nil {
			c.fatalc("Failed to install Docker engine: %s\n", err)
		}
		_ = c.startDockerEngine()
		version, err := c.getDockerServerVersion()
		if err != nil {
			c.fatalc("Failed to get Docker version: %s\n", err)
		}
		ok, err = c.validateInstalledDockerEngineVersion(version)
		if err != nil {
			c.fatalc("Failed to validate Docker version: %s\n", err)
		}
		if !ok {
			c.fatalc("Docker version does not meet the requirements: %s\n", version)
		}
	} else {
		c.printfc(colInfo, "Docker engine version is just fine\n")
	}
	if err := c.enableDockerEngine(); err != nil {
		c.fatalc("Failed to enable the Docker service: %s\n", err)
	}
	if c.isSwarmInit() {
		c.printfc(colInfo, "Swarm mode already set\n")
	} else {
		err = c.joinSwarm()
		if err != nil {
			c.fatalc("Failed to init or join the Swarm cluster: %s\n", err)
		}
	}
	if err := c.pullAmpImage(ampImageTag); err != nil {
		c.fatalc("%v", err)
	}
	if manager {
		c.clusterLoader.ampTag = ampImageTag
		if err := c.clusterLoader.init(c, ""); err != nil {
			c.fatalc("%v", err)
		}
		if err := c.clusterLoader.startClusterServices(); err != nil {
			c.fatalc("%v", err)
		}
	}
}
