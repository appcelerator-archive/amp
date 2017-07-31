package client

import (
	"io"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
)

func NewDockerCli(stdin io.ReadCloser, stdout, stderr io.Writer) *command.DockerCli {
	d := command.NewDockerCli(stdin, stdout, stderr)
	opts := flags.NewClientOptions()
	opts.Common = flags.NewCommonOptions()
	d.Initialize(opts)
	return d
}

