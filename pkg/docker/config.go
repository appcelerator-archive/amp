package docker

import (
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
	"golang.org/x/net/context"
)

func (d *Docker) ListConfigs() ([]swarm.Config, error) {
	return d.client.ConfigList(context.Background(), types.ConfigListOptions{})
}

func (d *Docker) ConfigExists(name string) (bool, error) {
	configs, err := d.ListConfigs()
	if err != nil {
		return false, err
	}
	for _, config := range configs {
		if config.Spec.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (d *Docker) CreateConfig(name string, data []byte) (string, error) {
	spec := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name: name,
		},
		Data: data,
	}
	config, err := d.client.ConfigCreate(context.Background(), spec)
	if err != nil {
		return "", err
	}
	return config.ID, nil
}
