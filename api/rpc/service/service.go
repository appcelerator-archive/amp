package service

import (
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var (
	// https://docs.docker.com/engine/reference/api/docker_remote_api/
	// `docker version` -> Server API version  => Docker 1.12x
	defaultVersion = "1.24"
	defaultHeaders = map[string]string{"User-Agent": "amplifier-1.0"}
	dockerSock     = "unix:///var/run/docker.sock"
	defaultNetwork = "amp-public"
	docker         *client.Client
	err            error
)

const serviceRoleLabelName = "io.amp.role"

// Service is used to implement ServiceServer
type Service struct{}

func init() {
	docker, err = client.NewClient(dockerSock, defaultVersion, nil, defaultHeaders)
	if err != nil {
		// fail fast
		log.Println("new client ....")
		panic(err)
	}
}

// Create implements ServiceServer
func (s *Service) Create(ctx context.Context, req *ServiceCreateRequest) (*ServiceCreateResponse, error) {
	log.Println(req)

	// TODO: pass-through right now, but will be refactored into a helper library
	response, err := CreateService(docker, ctx, req)
	return response, err
}

// CreateService uses docker api to create a service
func CreateService(docker *client.Client, ctx context.Context, req *ServiceCreateRequest) (*ServiceCreateResponse, error) {
	if req.ServiceSpec.Labels == nil {
		req.ServiceSpec.Labels = make(map[string]string)
	}
	req.ServiceSpec.Labels[serviceRoleLabelName] = "user"
	annotations := swarm.Annotations{
		Name:   req.ServiceSpec.Name,
		Labels: req.ServiceSpec.Labels,
	}

	containerSpec := swarm.ContainerSpec{
		Image: req.ServiceSpec.Image,
	}

	taskSpec := swarm.TaskSpec{
		ContainerSpec: containerSpec,
	}

	networks := []swarm.NetworkAttachmentConfig{
		{
			Target:  defaultNetwork,
			Aliases: []string{req.ServiceSpec.Name},
		},
	}

	mode := swarm.ServiceMode{
		Replicated: &swarm.ReplicatedService{
			Replicas: &req.ServiceSpec.Replicas,
		},
	}

	swarmSpec := swarm.ServiceSpec{
		Annotations:  annotations,
		TaskTemplate: taskSpec,
		Networks:     networks,
		Mode:         mode,
	}

	if req.ServiceSpec.PublishSpecs != nil {
		nn := len(req.ServiceSpec.PublishSpecs)
		if nn > 0 {
			swarmSpec.EndpointSpec = &swarm.EndpointSpec{
				Mode:  swarm.ResolutionModeVIP,
				Ports: make([]swarm.PortConfig, nn, nn),
			}
			for i, publish := range req.ServiceSpec.PublishSpecs {
				swarmSpec.EndpointSpec.Ports[i] = swarm.PortConfig{
					Name:          publish.Name,
					Protocol:      swarm.PortConfigProtocol(publish.Protocol),
					TargetPort:    publish.InternalPort,
					PublishedPort: publish.PublishPort,
				}
			}
		}
	}
	options := types.ServiceCreateOptions{}

	r, err := docker.ServiceCreate(ctx, swarmSpec, options)
	if err != nil {
		return nil, err
	}

	resp := &ServiceCreateResponse{
		Id: r.ID,
	}
	fmt.Printf("Service: %s created, id=%s\n", req.ServiceSpec.Name, resp.Id)
	return resp, nil
}

// Remove uses docker api to remove a service
func (s *Service) Remove(ctx context.Context, ID string) error {
	fmt.Printf("Service removed %s\n", ID)
	return docker.ServiceRemove(ctx, ID)
}
