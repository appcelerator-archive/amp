package stack

import (
	"fmt"

	"errors"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/compose/convert"
	"golang.org/x/net/context"
)

const (
	defaultNetworkDriver = "overlay"
)

type deployOptions struct {
	bundlefile       string
	composefile      string
	namespace        string
	sendRegistryAuth bool
	prune            bool
}

// NewDeployOptions {AMP} add constructor on this private struct
func NewDeployOptions(namespace string, composefile string, sendRegistryAuth bool) deployOptions {
	return deployOptions{
		composefile:      composefile,
		namespace:        namespace,
		sendRegistryAuth: sendRegistryAuth,
	}
}

// RunDeploy {AMP} make it public
func RunDeploy(dockerCli *command.DockerCli, opts deployOptions) error {
	ctx := context.Background()

	switch {
	case opts.bundlefile == "" && opts.composefile == "":
		return errors.New("Please specify either a bundle file (with --bundle-file) or a Compose file (with --compose-file).")
	case opts.bundlefile != "" && opts.composefile != "":
		return errors.New("You cannot specify both a bundle file and a Compose file.")
	case opts.bundlefile != "":
		return deployBundle(ctx, dockerCli, opts)
	default:
		return deployCompose(ctx, dockerCli, opts)
	}
}

// checkDaemonIsSwarmManager does an Info API call to verify that the daemon is
// a swarm manager. This is necessary because we must create networks before we
// create services, but the API call for creating a network does not return a
// proper status code when it can't create a network in the "global" scope.
func checkDaemonIsSwarmManager(ctx context.Context, dockerCli *command.DockerCli) error {
	info, err := dockerCli.Client().Info(ctx)
	if err != nil {
		return err
	}
	if !info.Swarm.ControlAvailable {
		return errors.New("This node is not a swarm manager. Use \"docker swarm init\" or \"docker swarm join\" to connect this node to swarm and try again.")
	}
	return nil
}

// pruneServices removes services that are no longer referenced in the source
func pruneServices(ctx context.Context, dockerCli command.Cli, namespace convert.Namespace, services map[string]struct{}) bool {
	client := dockerCli.Client()

	oldServices, err := getServices(ctx, client, namespace.Name())
	if err != nil {
		fmt.Fprintf(dockerCli.Err(), "Failed to list services: %s", err)
		return true
	}

	pruneServices := []swarm.Service{}
	for _, service := range oldServices {
		if _, exists := services[namespace.Descope(service.Spec.Name)]; !exists {
			pruneServices = append(pruneServices, service)
		}
	}
	return removeServices(ctx, dockerCli, pruneServices)
}
