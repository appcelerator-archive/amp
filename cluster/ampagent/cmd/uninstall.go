package cmd

import (
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/go-docker/cli/cli/command/stack"
	"github.com/appcelerator/go-docker/docker/pkg/term"
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
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := docker.NewDockerCli(stdin, stdout, stderr)

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

	if err := removeVolumes(); err != nil {
		return err
	}

	return removeInitialNetworks()
}
