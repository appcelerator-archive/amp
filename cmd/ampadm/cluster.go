package main

import (
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

const (
	ampVersion          = "0.5.0"
	infraPrivateNetwork = "amp-infra"
	infraPublicNetwork  = "amp-public"
)

type clusterStack struct {
	serviceMap map[string]*clusterService
	imageMap   map[string]string
	networks   []string
}

type clusterService struct {
	id              string
	name            string
	image           string
	desiredReplicas int
	spec            *swarm.ServiceSpec
}

func (s *clusterStack) addService(name string, imageLabel string, replicas int, spec *swarm.ServiceSpec) {
	image, ok := s.imageMap[imageLabel]
	if !ok {
		fmt.Printf("The image doesn't exist for service: %s", name)
		return
	}
	spec.Annotations.Name = name
	spec.TaskTemplate.ContainerSpec.Image = image
	if replicas > 0 {
		nb := uint64(replicas)
		spec.Mode = swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &nb,
			},
		}
	} else {
		spec.Mode = swarm.ServiceMode{
			Global: &swarm.GlobalService{},
		}
	}
	labels := spec.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	s.serviceMap[name] = &clusterService{
		name:            name,
		image:           image,
		desiredReplicas: replicas,
		spec:            spec,
	}
}

func (s *clusterStack) init(local bool, ampTag string) {
	s.serviceMap = make(map[string]*clusterService)
	s.imageMap = make(map[string]string)
	s.networks = []string{infraPrivateNetwork, infraPublicNetwork}
	if local {
		s.imageMap["amp"] = "appcelerator/amp:local"
	} else if ampTag != "" {
		s.imageMap["amp"] = fmt.Sprintf("appcelerator/amp:%s", ampTag)
	} else {
		s.imageMap["amp"] = fmt.Sprintf("appcelerator/amp:%s", ampVersion)
	}
	s.addService("adm-server", "amp", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
				//"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{"adm-server"},
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeBind,
							Source: "/var/run/docker.sock",
							Target: "/var/run/docker.sock",
						},
					},
				},
				Placement: &swarm.Placement{
					Constraints: []string{"node.role == manager"},
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Mode: swarm.ResolutionModeVIP,
				Ports: []swarm.PortConfig{
					{
						TargetPort:    31315,
						PublishedPort: 31315,
					},
				},
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"adm-server"},
				},
				{
					Target:  infraPublicNetwork,
					Aliases: []string{"adm-server"},
				},
			},
		})

	s.addService("adm-agent", "amp", 0,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
				//"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{"adm-agent"},
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeBind,
							Source: "/var/run/docker.sock",
							Target: "/var/run/docker.sock",
						},
					},
				},
				Placement: nil,
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"adm-agent"},
				},
			},
		})
}
