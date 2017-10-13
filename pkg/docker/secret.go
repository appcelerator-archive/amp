package docker

import (
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
	"golang.org/x/net/context"
)

func (d *Docker) ListSecrets() ([]swarm.Secret, error) {
	return d.client.SecretList(context.Background(), types.SecretListOptions{})
}

func (d *Docker) SecretExists(name string) (bool, error) {
	secrets, err := d.ListSecrets()
	if err != nil {
		return false, err
	}
	for _, secret := range secrets {
		if secret.Spec.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (d *Docker) CreateSecret(name string, data []byte) (string, error) {
	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name: name,
		},
		Data: data,
	}
	secret, err := d.client.SecretCreate(context.Background(), spec)
	if err != nil {
		return "", err
	}
	return secret.ID, nil
}
