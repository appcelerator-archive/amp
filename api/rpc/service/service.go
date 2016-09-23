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

	serv := req.ServiceSpec
	//prepare swarm.ServiceSpec full instance
	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   serv.Name,
			Labels: make(map[string]string),
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image:           serv.Image,
				Args:            nil, //[]string
				Env:             nil, //[]string
				Labels:          serv.ContainerLabels,
				Dir:             "",
				User:            "",
				Groups:          nil, //[]string
				Mounts:          nil, //[]mount.Mount
				StopGracePeriod: nil, //*time.Duration
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  defaultNetwork,
					Aliases: []string{req.ServiceSpec.Name},
				},
			},
			Resources:     nil, //*ResourceRequirements
			RestartPolicy: nil, //*RestartPolicy
			Placement: &swarm.Placement{
				Constraints: nil, //[]string
			},
			LogDriver: nil, //*Driver
		},
		Networks: nil, //[]NetworkAttachmentConfig
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &serv.Replicas,
			},
		},
		UpdateConfig: &swarm.UpdateConfig{
			Parallelism:   0,
			Delay:         0,
			FailureAction: "",
		},
		EndpointSpec: nil, // &EndpointSpec
	}
	//add common labels
	if service.Annotations.Labels == nil {
		service.Annotations.Labels = make(map[string]string)
	}
	service.Annotations.Labels[serviceRoleLabelName] = "user"

	if req.ServiceSpec.PublishSpecs != nil {
		nn := len(req.ServiceSpec.PublishSpecs)
		if nn > 0 {
			service.EndpointSpec = &swarm.EndpointSpec{
				Mode:  swarm.ResolutionModeVIP,
				Ports: make([]swarm.PortConfig, nn, nn),
			}
			for i, publish := range req.ServiceSpec.PublishSpecs {
				service.EndpointSpec.Ports[i] = swarm.PortConfig{
					Name:          publish.Name,
					Protocol:      swarm.PortConfigProtocol(publish.Protocol),
					TargetPort:    publish.InternalPort,
					PublishedPort: publish.PublishPort,
				}
			}
		}
	}
	options := types.ServiceCreateOptions{}

	r, err := docker.ServiceCreate(ctx, service, options)
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
