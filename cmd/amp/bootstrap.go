package main

import (
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var bootstrapCmd = &cobra.Command{
	Use:     "bootstrap",
	Short:   "Bootstraps the node",
	Long:    `Installs Docker, setups Swarm, run the AMP agent`,
	Aliases: []string{"bs"},
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap(cmd, args)
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
// if on a manager, launch the stack ampcore
func bootstrap(cmd *cobra.Command, args []string) {
	c := newManager(cmd.Flag("verbose").Value.String())
	c.printf(colInfo, "Bootstrapping the node with CLI version %s\n", Version)

	var version string

	_ = c.startDockerEngine()
	version, _ = c.getDockerServerVersion()
	ok, err := c.validateInstalledDockerEngineVersion(version)
	if err != nil {
		c.fatalf("Failed to validate Docker version: %s\n", err)
	}
	if !ok {
		if version != "0.0.0" {
			c.printf(colInfo, "Installed version %s is lower than %s\n", version, minimumDockerEngineVersion)
		}
		if err := c.installDockerEngine(); err != nil {
			c.fatalf("Failed to install Docker engine: %s\n", err)
		}
		_ = c.startDockerEngine()
		version, err := c.getDockerServerVersion()
		if err != nil {
			c.fatalf("Failed to get Docker version: %s\n", err)
		}
		ok, err = c.validateInstalledDockerEngineVersion(version)
		if err != nil {
			c.fatalf("Failed to validate Docker version: %s\n", err)
		}
		if !ok {
			c.fatalf("Docker version does not meet the requirements: %s\n", version)
		}
	} else {
		c.printf(colInfo, "Docker engine version is just fine\n")
	}
	if err := c.enableDockerEngine(); err != nil {
		c.fatalf("Failed to enable the Docker service: %s\n", err)
	}
	if c.isSwarmInit() {
		c.printf(colInfo, "Swarm mode already set\n")
	} else {
		err = c.joinSwarm()
		if err != nil {
			c.fatalf("Failed to init or join the Swarm cluster: %s\n", err)
		}
	}
	if err := c.pullAmpImage(ampImageTag); err != nil {
		c.fatalf("%v", err)
	}
	if manager {
		smanager := newManager("")
		if err := smanager.connectDocker(); err != nil {
			smanager.fatalf("Needs AMP installed to start ampcore")
		}
		if err := smanager.systemPrerequisites(); err != nil {
			smanager.fatalf("Prerequiste error: %v\n", err)
		}
		smanager.printf(colRegular, "starting stack: ampcore\n")
		data, err := stack.ResolvedComposeFileVariables("ampcore.yml", "amp.var", "")
		if err != nil {
			smanager.fatalf("start ampCore error: %v\n", err)
		}
		stackInstance := &stack.Stack{
			Name:     "ampcore",
			FileData: data,
		}
		server := stack.NewServer(nil, smanager.docker)
		request := &stack.StackDeployRequest{Stack: stackInstance, RegistryAuth: false}
		reply, err := server.Deploy(context.Background(), request)
		if err != nil {
			smanager.fatalf("start ampCore error: %v\n", err)
		}
		smanager.printf(colSuccess, reply.Answer)
	}
}
