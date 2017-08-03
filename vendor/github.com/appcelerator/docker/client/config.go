package client

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

type configCreateOptions struct {
	name   string
	file   string
	labels opts.ListOpts
}

// ConfigCreate is intended to be used as the client from the amplifier API
func ConfigCreate(dockerCli command.Cli, name string, labels map[string]string, data []byte) (string, error) {
	cli := dockerCli.Client()
	ctx := context.Background()

	spec := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name:   name,
			Labels: labels,
		},
		Data: data,
	}
	resp, err := cli.ConfigCreate(ctx, spec)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// ConfigList is intended to be used as the client from the amplifier API
func ConfigList(dockerCli command.Cli) ([]string, error) {
	cli := dockerCli.Client()
	ctx := context.Background()

	opts := opts.NewFilterOpt()

	list, err := cli.ConfigList(ctx, types.ConfigListOptions{Filters: opts.Value()})
	if err != nil {
		return nil, err
	}

	secrets := []string{}
	for _, secret := range list {
		secrets = append(secrets, secret.Spec.Name)
	}
	return secrets, nil
}

// ConfigRemove is intended to be used as the client from the amplifier API
func ConfigRemove(dockerCli command.Cli, id string) error {
	cli := dockerCli.Client()
	ctx := context.Background()
	return cli.ConfigRemove(ctx, id)
}
