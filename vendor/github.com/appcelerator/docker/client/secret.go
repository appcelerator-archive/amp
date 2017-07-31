package client

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

type createOptions struct {
	name   string
	file   string
	labels opts.ListOpts
}

// SecretCreate is intended to be used as the client from the amplifier API
func SecretCreate(dockerCli command.Cli, name string, labels map[string]string, data []byte) (string, error) {
	cli := dockerCli.Client()
	ctx := context.Background()


	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name: name,
			Labels: labels,
		},
		Data: data,
	}
	resp, err := cli.SecretCreate(ctx, spec)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
