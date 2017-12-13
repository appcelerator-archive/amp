package cmd

import (
	"github.com/appcelerator/amp/docker/cli/cli/command/stack"
	"github.com/appcelerator/amp/docker/docker/pkg/term"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/spf13/cobra"
)

func NewUninstallCommand() *cobra.Command {
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall amp services from swarm environment",
		RunE:  Uninstall,
	}
	return uninstallCmd
}

func Uninstall(cmd *cobra.Command, args []string) error {
	Docker = docker.NewEnvClient()
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := docker.NewDockerCli(stdin, stdout, stderr)
	if err := Docker.Connect(); err != nil {
		return err
	}

	namespace := "amp"
	if len(args) > 0 {
		namespace = args[0]
	}

	opts := stack.RemoveOptions{
		Namespaces: []string{namespace},
	}

	if err := stack.RunRemove(dockerCli, opts); err != nil {
		return err
	}

	// workaround for https://github.com/moby/moby/issues/32620
	if err := removeExitedContainers(30); err != nil {
		return err
	}

	if err := removeVolumes(5); err != nil {
		return err
	}

	return removeInitialNetworks()
}
