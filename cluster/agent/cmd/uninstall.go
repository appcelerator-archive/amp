package main

import (
	"github.com/appcelerator/amp/cluster/agent/pkg/docker"
	"github.com/appcelerator/amp/cluster/agent/pkg/docker/stack"
	"github.com/docker/docker/pkg/term"
	"github.com/spf13/cobra"
)

func NewUninstallCommand() *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall amp services from swarm environment",
		RunE:  uninstall,
	}
	return uninstallCmd
}

func uninstall(cmd *cobra.Command, args []string) error {
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := docker.NewDockerCli(stdin, stdout, stderr)

	namespace := "amp"
	if len(args) > 0 {
		namespace = args[0]
	}

	opts := stack.RemoveOptions{
		Namespaces: []string{namespace},
	}

	if err := stack.Remove(dockerCli, opts); err != nil {
		return err
	}

	if err := removeVolumes(); err != nil {
		return err
	}

	return removeInitialNetworks()
}
