package servercore

import (
	"github.com/appcelerator/amp/config"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

const (
	infraPrivateNetwork = "amp-infra"
	infraPublicNetwork  = "amp-public"
)

const (
	ampVersion           = "0.5.0"
	uiVersion            = "0.3.0"
	influxDBVersion      = "1.1.2"
	kapacitorVersion     = "1.1.2"
	telegrafVersion      = "1.1.1"
	grafanaVersion       = "1.1.0"
	elasticsearchVersion = "5.1.1"
	etcdVersion          = "3.1.0-rc.1"
	natsVersion          = "0.3.0"
	haproxyVersion       = "1.0.3"
	registryVersion      = "2.5.1"
)

func getAMPInfrastructureStack(m *AMPInfraManager) *ampStack {
	//init stack
	stack := ampStack{}
	stack.init()

	//add images
	if m.Local {
		stack.addImage("amp", "appcelerator/amp:local")
	} else {
		stack.addImage("amp", "appcelerator/amp:"+ampVersion)
	}
	stack.addImage("amp-ui", "appcelerator/amp-ui:"+uiVersion)
	stack.addImage("elasticsearch", "appcelerator/elasticsearch-amp:"+elasticsearchVersion)
	stack.addImage("grafana", "appcelerator/grafana-amp:"+grafanaVersion)
	stack.addImage("haproxy", "appcelerator/haproxy:"+haproxyVersion)
	stack.addImage("influxdb", "appcelerator/influxdb-amp:"+influxDBVersion)
	stack.addImage("kapacitor", "appcelerator/kapacitor-amp:"+kapacitorVersion)
	stack.addImage("telegraf", "appcelerator/telegraf:telegraf-"+telegrafVersion)
	stack.addImage("registry", "registry:"+registryVersion)
	stack.addImage("etcd", "appcelerator/etcd:"+etcdVersion)
	stack.addImage("pinger", "appcelerator/pinger:latest")
	stack.addImage("nats", "appcelerator/amp-nats-streaming:"+natsVersion)

	//add networks
	stack.networks = []string{infraPrivateNetwork, infraPublicNetwork}

	//add volumes to clean
	stack.volumesToRemove = []string{"amp-etcd"}

	//add etcd
	stack.addService(m, "etcd", "etcd", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{
						"--name etcd",
						"--listen-client-urls http://0.0.0.0:2379",
						"--advertise-client-urls " + amp.EtcdDefaultEndpoint,
					},
					Env: nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeVolume,
							Source: "amp-etcd",
							Target: "/data",
						},
					},
				},
				Placement: &swarm.Placement{
					Constraints: nil,
				},
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"etcd"},
				},
			},
		})

	//add haproxy
	stack.addService(m, "haproxy", "haproxy", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: &swarm.Placement{
					Constraints: nil,
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Mode: swarm.ResolutionModeVIP,
				Ports: []swarm.PortConfig{
					{
						TargetPort:    8080,
						PublishedPort: 8080,
					},
					{
						TargetPort:    80,
						PublishedPort: 80,
					},
				},
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"haproxy"},
				},
				{
					Target:  infraPublicNetwork,
					Aliases: []string{"haproxy"},
				},
			},
		},
		"etcd")

	//add amp-ui
	stack.addService(m, "amp-ui", "amp-ui", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
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
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"amp-ui"},
				},
			},
		})

	//add nats
	stack.addService(m, "nats", "nats", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: &swarm.Placement{
					Constraints: nil,
				},
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"nats"},
				},
			},
		})

	//add influxdb
	stack.addService(m, "influxdb", "influxdb", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: &swarm.Placement{
					Constraints: nil,
				},
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"influxdb"},
				},
			},
		})

	//add elasticsearch
	stack.addService(m, "elasticsearch", "elasticsearch", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: &swarm.Placement{
					Constraints: nil,
				},
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"elasticsearch"},
				},
			},
		})

	//add amp-agent
	stack.addService(m, "amp-agent", "amp", 0,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{"amp-agent"},
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
					Aliases: []string{"amp-agent"},
				},
			},
		},
		"nats", "elasticsearch")

	//add registry
	stack.addService(m, "registry", "registry", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeVolume,
							Source: "amp-registry",
							Target: "/var/lib/registry",
						},
					},
				},
				Placement: nil,
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"registry"},
				},
				{
					Target:  infraPublicNetwork,
					Aliases: []string{"registry"},
				},
			},
		})

	//add amp-log-worker
	stack.addService(m, "amp-log-worker", "amp", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{"amp-log-worker"},
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: nil,
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"amp-log-worker"},
				},
			},
		},
		"nats", "elasticsearch")

	//add grafana
	stack.addService(m, "grafana", "grafana", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: nil,
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"grafana"},
				},
			},
		},
		"influxdb")

	//add kapacitor
	stack.addService(m, "kapacitor", "kapacitor", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: nil,
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"kapacitor"},
				},
			},
		},
		"influxdb")

	//add telegraf-agent
	stack.addService(m, "telegraf-agent", "telegraf", 0,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env: []string{
						"OUTPUT_INFLUXDB_ENABLED=true",
						"INFLUXDB_URL=" + amp.InfluxDefaultURL,
						"TAG_datacenter=dc1",
						"TAG_type=core",
						"INPUT_DOCKER_ENABLED=true",
						"INPUT_CPU_ENABLED=true",
						"INPUT_DISK_ENABLED=true",
						"INPUT_DISKIO_ENABLED=true",
						"INPUT_KERNEL_ENABLED=true",
						"INPUT_MEM_ENABLED=true",
						"INPUT_PROCESS_ENABLED=true",
						"INPUT_SWAP_ENABLED=true",
						"INPUT_SYSTEM_ENABLED=true",
						"INPUT_NET_ENABLED=true",
						"INPUT_HAPROXY_ENABLED=false",
						"INFLUXDB_TIMEOUT=20",
					},
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: []mount.Mount{
						{
							Type:   mount.TypeBind,
							Source: "/var/run/docker.sock",
							Target: "/var/run/docker.sock",
						},
						{
							Type:   mount.TypeBind,
							Source: "/var/run/utmp",
							Target: "/var/run/utmp",
						},
					},
				},
				Placement: nil,
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"amp-agent"},
				},
			},
		}, "influxdb")

	//add telegraf-haproxy
	stack.addService(m, "telegraf-haproxy", "telegraf", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
					Env: []string{
						"OUTPUT_INFLUXDB_ENABLED=true",
						"INFLUXDB_URL=" + amp.InfluxDefaultURL,
						"INPUT_DOCKER_ENABLED=false",
						"INPUT_CPU_ENABLED=false",
						"INPUT_NET_ENABLED=false",
						"INPUT_DISK_ENABLED=false",
						"INPUT_DISKIO_ENABLED=false",
						"INPUT_KERNEL_ENABLED=false",
						"INPUT_MEM_ENABLED=false",
						"INPUT_PROCESS_ENABLED=false",
						"INPUT_SWAP_ENABLED=false",
						"INPUT_SYSTEM_ENABLED=false",
						"INPUT_HAPROXY_ENABLED=true",
						"INPUT_HAPROXY_SERVER=http://haproxy:8082/admin?stats",
						"INFLUXDB_TIMEOUT=20",
					},
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: nil,
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"amp-agent"},
				},
			},
		},
		"haproxy", "influxdb")

	//add amplifier
	stack.addService(m, "amplifier", "amp", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: nil,
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
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"amplifier"},
				},
				{
					Target:  infraPublicNetwork,
					Aliases: []string{"amplifier"},
				},
			},
		},
		"etcd", "nats", "elasticsearch", "influxdb")

	//add amplifier-gateway
	stack.addService(m, "amplifier-gateway", "amp", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{
						"amplifier-gateway",
						"--amplifier_endpoint",
						amp.AmplifierDefaultEndpoint,
					},
					Env: nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
				},
			},
			EndpointSpec: nil,
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPublicNetwork,
					Aliases: []string{"amplifier-gateway"},
				},
			},
		},
		"amplifier")

	//add amp-function-listener
	stack.addService(m, "amp-function-listener", "amp", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{"amp-function-listener"},
					Env:  nil,
					Labels: map[string]string{
						"io.amp.role": "infrastructure",
					},
					Mounts: nil,
				},
				Placement: nil,
			},
			EndpointSpec: &swarm.EndpointSpec{
				Mode: swarm.ResolutionModeVIP,
				Ports: []swarm.PortConfig{
					{
						TargetPort:    80,
						PublishedPort: 4242,
					},
				},
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target:  infraPrivateNetwork,
					Aliases: []string{"amp-function-listener"},
				},
			},
		},
		"nats", "etcd")

	//add amp-function-worker
	stack.addService(m, "amp-function-worker", "amp", 1,
		&swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Labels: map[string]string{
					"io.amp.role": "infrastructure",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Args: []string{"amp-function-worker"},
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
					Aliases: []string{"amp-function-worker"},
				},
			},
		},
		"nats")

	//return stack
	return &stack
}
