package service

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var (
	defaultNetwork = "amp-public"
	err            error
)

const serviceRoleLabelName = "io.amp.role"

// Service is used to implement ServiceServer
type Service struct {
	Docker *client.Client
}

// SwarmMode is needed to export isServiceSpec_Mode type, which consumers can use to
// create a variable and assign either a ServiceSpec_Replicated or ServiceSpec_Global struct
type SwarmMode isServiceSpec_Mode

func init() {
	//Nothing to do for now
}

// Create uses docker api to create a service
func (s *Service) Create(ctx context.Context, req *ServiceCreateRequest) (*ServiceCreateResponse, error) {

	serv := req.ServiceSpec

	var serviceMode swarm.ServiceMode
	switch mode := serv.Mode.(type) {
	case *ServiceSpec_Replicated:
		serviceMode = swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &mode.Replicated.Replicas,
			},
		}
	case *ServiceSpec_Global:
		serviceMode = swarm.ServiceMode{
			Global: &swarm.GlobalService{},
		}
	}

	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   serv.Name,
			Labels: serv.Labels,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image:           serv.Image,
				Args:            nil, //[]string
				Env:             serv.Env,
				Labels:          serv.ContainerLabels,
				Dir:             "",
				User:            "",
				Groups:          nil, //[]string
				Mounts:          nil, //[]mount.Mount
				StopGracePeriod: nil, //*time.Duration
			},
			Networks: nil,
			Resources:     nil, //*ResourceRequirements
			RestartPolicy: nil, //*RestartPolicy
			Placement: &swarm.Placement{
				Constraints: nil, //[]string
			},
			LogDriver: nil, //*Driver
		},
		Networks: []swarm.NetworkAttachmentConfig{
			{
				Target:  defaultNetwork,
				Aliases: []string{req.ServiceSpec.Name},
			},
		},
		UpdateConfig: &swarm.UpdateConfig{
			Parallelism:   0,
			Delay:         0,
			FailureAction: "",
		},
		EndpointSpec: nil, // &EndpointSpec
		Mode:         serviceMode,
	}

	// add network
	if req.ServiceSpec.Networks != nil {
		service.Networks = make([]swarm.NetworkAttachmentConfig, len(req.ServiceSpec.Networks), len(req.ServiceSpec.Networks))
		for i, net := range(req.ServiceSpec.Networks) {
			fmt.Printf("network: %v\n", net)
			service.Networks[i] = swarm.NetworkAttachmentConfig{
				Target:  net.Target,
				Aliases: net.Aliases,
			}
		}
	}

	// ensure supplied service label map is not nil, then add custom amp labels
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

	r, err := s.Docker.ServiceCreate(ctx, service, options)
	if err != nil {
		return nil, err
	}

	resp := &ServiceCreateResponse{
		Id: r.ID,
	}
	fmt.Printf("Service: %s created, id=%s\n", req.ServiceSpec.Name, resp.Id)
	return resp, nil
}

// Remove implements ServiceServer
func (s *Service) Remove(ctx context.Context, req *RemoveRequest) (*RemoveResponse, error) {
	err := s.Docker.ServiceRemove(ctx, req.Ident)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Service removed %s\n", req.Ident)
	response := &RemoveResponse{
		Ident: req.Ident,
	}

	return response, nil
}
