package stack

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/command"
	"golang.org/x/net/context"
)

type removeOptions struct {
	namespaces []string
}

// NewRemoveOptions {AMP}: add constructor for this private struct
func NewRemoveOptions(namespace []string) removeOptions {
	return removeOptions{
		namespaces: namespace,
	}
}

// RunRemove {AMP}: make it public
func RunRemove(dockerCli command.Cli, opts removeOptions) error {
	namespaces := opts.namespaces
	client := dockerCli.Client()
	ctx := context.Background()

	var errs []string
	for _, namespace := range namespaces {
		services, err := getServices(ctx, client, namespace)
		if err != nil {
			return err
		}

		networks, err := getStackNetworks(ctx, client, namespace)
		if err != nil {
			return err
		}

		secrets, err := getStackSecrets(ctx, client, namespace)
		if err != nil {
			return err
		}

		if len(services)+len(networks)+len(secrets) == 0 {
			fmt.Fprintf(dockerCli.Out(), "Nothing found in stack: %s\n", namespace)
			continue
		}

		hasError := removeServices(ctx, dockerCli, services)
		hasError = removeSecrets(ctx, dockerCli, secrets) || hasError
		hasError = removeNetworks(ctx, dockerCli, networks) || hasError

		if hasError {
			errs = append(errs, fmt.Sprintf("Failed to remove some resources from stack: %s", namespace))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}

func removeServices(
	ctx context.Context,
	dockerCli command.Cli,
	services []swarm.Service,
) bool {
	var err error
	for _, service := range services {
		fmt.Fprintf(dockerCli.Err(), "Removing service %s\n", service.Spec.Name)
		if err = dockerCli.Client().ServiceRemove(ctx, service.ID); err != nil {
			fmt.Fprintf(dockerCli.Err(), "Failed to remove service %s: %s", service.ID, err)
		}
	}
	return err != nil
}

func removeNetworks(
	ctx context.Context,
	dockerCli command.Cli,
	networks []types.NetworkResource,
) bool {
	var err error
	for _, network := range networks {
		fmt.Fprintf(dockerCli.Err(), "Removing network %s\n", network.Name)
		if err = dockerCli.Client().NetworkRemove(ctx, network.ID); err != nil {
			fmt.Fprintf(dockerCli.Err(), "Failed to remove network %s: %s", network.ID, err)
		}
	}
	return err != nil
}

func removeSecrets(
	ctx context.Context,
	dockerCli command.Cli,
	secrets []swarm.Secret,
) bool {
	var err error
	for _, secret := range secrets {
		fmt.Fprintf(dockerCli.Err(), "Removing secret %s\n", secret.Spec.Name)
		if err = dockerCli.Client().SecretRemove(ctx, secret.ID); err != nil {
			fmt.Fprintf(dockerCli.Err(), "Failed to remove secret %s: %s", secret.ID, err)
		}
	}
	return err != nil
}
