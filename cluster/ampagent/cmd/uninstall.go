package cmd

import (
	"context"
	"log"
	"strings"
	"time"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"github.com/appcelerator/amp/docker/cli/cli/command/stack"
	"github.com/appcelerator/amp/docker/cli/opts"
	"github.com/appcelerator/amp/docker/docker/pkg/term"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/spf13/cobra"
)

const ampVolumesPrefix = "amp_"

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

	// workaround for https://github.com/moby/moby/issues/32620
	if err := removeExitedContainers(30); err != nil {
		return err
	}

	if err := removeVolumes(5); err != nil {
		return err
	}

	return removeInitialNetworks()
}

func removeExitedContainers(timeout int) error {
	i := 0
	dontKill := []string{"amp-agent", "amp-local"}
	var containers []types.Container
	if timeout == 0 {
		timeout = 30 // default value
	}
	log.Println("waiting for all services to clear up...")
	filter := filters.NewArgs()
	filter.Add("is-task", "true")
	filter.Add("label", "io.amp.role=infrastructure")
	for i < timeout {
		containers, err := Docker.GetClient().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return err
		}
		if len(containers) == 0 {
			log.Println("cleared up")
			break
		}
		for _, c := range containers {
			switch c.State {
			case "exited":
				log.Printf("Removing container %s [%s]\n", c.Names[0], c.Status)
				err := Docker.GetClient().ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{})
				if err != nil {
					if strings.Contains(err.Error(), "already in progress") {
						continue // leave it to Docker
					}
					return err
				}
			case "removing", "running":
				// ignore it, _running_ containers will be killed after the loop
				// _removing_ containers are in progress of deletion
			default:
				// this is not expected
				log.Printf("Container %s found in status %s, %s\n", c.Names[0], c.Status, c.State)
			}
		}
		i++
		time.Sleep(1 * time.Second)
	}
	containers, err := Docker.GetClient().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return err
	}
	if i == timeout {
		log.Println("timing out")
		log.Printf("%d containers left\n", len(containers))
	}
	//
	for _, c := range containers {
		for _, e := range dontKill {
			if strings.Contains(c.Names[0], e) {
				continue
			}
		}
		log.Printf("Force removing container %s [%s]", c.Names[0], c.State)
		if err := Docker.GetClient().ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			if strings.Contains(err.Error(), "already in progress") {
				continue // leave it to Docker
			}
			return err
		}
	}
	return nil
}

func removeInitialNetworks() error {
	for _, network := range ampnetworks {
		// Check if network already exists
		id, err := Docker.NetworkID(network)
		if err != nil {
			return err
		}
		if id == "" {
			continue // Skipping non existent network
		}

		// Remove network
		if err := Docker.RemoveNetwork(id); err != nil {
			return err
		}
		log.Printf("Successfully removed network %s [%s]", network, id)
	}
	return nil
}

func removeVolumes(timeout int) error {
	// volume remove timeout (sec)
	if timeout == 0 {
		timeout = 5 // default value
	}
	// List amp volumes
	filter := opts.NewFilterOpt()
	filter.Set("name=" + ampVolumesPrefix)
	volumes, err := Docker.ListVolumes(filter)
	if err != nil {
		return nil
	}
	// Remove volumes
	for _, volume := range volumes {
		log.Printf("Removing volume [%s]... ", volume.Name)
		if err := Docker.RemoveVolume(volume.Name, false, timeout); err != nil {
			log.Println("Failed")
			return err
		}
	}
	return nil
}
