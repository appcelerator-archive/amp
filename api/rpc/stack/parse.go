package stack

import (
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/rpc/service"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

type stackSpec struct {
	Services map[string]serviceSpec `yaml:"services"`
	Networks map[string]networkSpec `yaml:"networks"`
}

type serviceSpec struct {
	Image           string                    `yaml:"image"`
	Public          []publishSpec             `yaml:"public"`
	Mode            string                    `yaml:"mode"`
	Replicas        uint64                    `yaml:"replicas"`
	Args            interface{}               `yaml:"args"`
	Environment     interface{}               `yaml:"environment"`
	Labels          interface{}               `yaml:"labels"`
	ContainerLabels interface{}               `yaml:"container_labels"`
	Networks        map[string]networkAliases `yaml:"networks"`
	Mounts          []string                  `yaml:"volumes"`
}

type publishSpec struct {
	Name         string `yaml:"name"`
	Protocol     string `yaml:"protocol"`
	PublishPort  uint32 `yaml:"publish_port"`
	InternalPort uint32 `yaml:"internal_port"`
}

type networkAliases struct {
	Aliases []string `yaml:"aliases"`
}

type networkSpec struct {
	External   interface{}       `yaml:"external"`
	Driver     string            `yaml:"driver"`
	EnableIPv6 bool              `yaml:"enable_ipv6"`
	IPAM       *networkIPAM      `yaml:"ipam"`
	Internal   bool              `yaml:"internal"`
	Options    map[string]string `yaml:"driver_opts"`
	Labels     map[string]string `yaml:"labels"`
}

// IPAM represents IP Address Management
type networkIPAM struct {
	Driver  string            `yaml:"driver"`
	Options map[string]string `yaml:"options"`
	Config  []ipamConfig      `yaml:"config"`
}

// IPAMConfig represents IPAM configurations
type ipamConfig struct {
	Subnet     string            `yaml:"subnet"`
	IPRange    string            `yaml:"ip_range"`
	Gateway    string            `yaml:"gateway"`
	AuxAddress map[string]string `yaml:"aux_address"`
}

// ParseStackfile main function to parse stackfile
func ParseStackfile(ctx context.Context, in string) (*Stack, error) {
	var stack = &Stack{}
	specs, err := parseStack([]byte(in))
	if err != nil {
		return nil, err
	}
	networkMap, err := copyNetworks(stack, specs.Networks)
	if err != nil {
		return nil, err
	}
	if err := copyServices(stack, specs.Services, networkMap); err != nil {
		return nil, err
	}
	return stack, err
}

func parseStack(b []byte) (*stackSpec, error) {
	var specs stackSpec
	if err := yaml.Unmarshal(b, &specs); err != nil {
		return nil, err
	}
	return &specs, nil
}

func copyNetworks(stack *Stack, specs map[string]networkSpec) (map[string]string, error) {
	networkMap := make(map[string]string)
	for name, spec := range specs {
		external := "false"
		if extMap, ok := spec.External.(map[interface{}]interface{}); ok {
			external = extMap["name"].(string)
			networkMap[name] = external
		} else if ext, ok := spec.External.(bool); ok {
			external = fmt.Sprintf("%t", ext)
		} else if spec.External != nil {
			return networkMap, fmt.Errorf("invalid syntax near networks: %s: external", name)
		}
		stack.Networks = append(stack.Networks, &NetworkSpec{
			External:   external,
			Name:       name,
			Driver:     spec.Driver,
			EnableIpv6: spec.EnableIPv6,
			Ipam:       copyIPAM(spec.IPAM),
			Internal:   spec.Internal,
			Options:    spec.Options,
			Labels:     spec.Labels,
		})
	}
	return networkMap, nil

}

func copyIPAM(ipam *networkIPAM) *NetworkIPAM {
	if ipam == nil {
		return nil
	}
	return &NetworkIPAM{
		Driver:  ipam.Driver,
		Options: ipam.Options,
		Config:  copyIPAMConfig(ipam.Config),
	}
}

func copyIPAMConfig(config []ipamConfig) []*NetworkIPAMConfig {
	configList := []*NetworkIPAMConfig{}
	if config != nil {
		for _, conf := range config {
			configList = append(configList, &NetworkIPAMConfig{
				Subnet:     conf.Subnet,
				IpRange:    conf.IPRange,
				Gateway:    conf.Gateway,
				AuxAddress: conf.AuxAddress,
			})
		}
	}
	return configList
}

func copyServices(stack *Stack, specs map[string]serviceSpec, networkMap map[string]string) error {
	for name, spec := range specs {
		// try to parse arguments entries as a map
		// else try to parse environment as string entries
		args := []string{}
		if argMap, ok := spec.Args.(map[interface{}]interface{}); ok {
			for k, v := range argMap {
				args = append(args, k.(string)+"="+v.(string))
			}
		} else if argList, ok := spec.Args.([]interface{}); ok {
			for _, e := range argList {
				args = append(args, e.(string))
			}
		}

		// try to parse environment entries as a map
		// else try to parse environment as string entries
		env := []string{}
		if envMap, ok := spec.Environment.(map[interface{}]interface{}); ok {
			for k, v := range envMap {
				env = append(env, k.(string)+"="+v.(string))
			}
		} else if envList, ok := spec.Environment.([]interface{}); ok {
			for _, e := range envList {
				env = append(env, e.(string))
			}
		}
		// try to parse labels as a map
		// else try to parse labels as string entries
		var labels = map[string]string{}
		if labelMap, ok := spec.Labels.(map[interface{}]interface{}); ok {
			for k, v := range labelMap {
				labels[k.(string)] = v.(string)
			}
		} else if labelList, ok := spec.Labels.([]interface{}); ok {
			for _, s := range labelList {
				a := strings.Split(s.(string), "=")
				labels[a[0]] = a[1]
			}
		}
		// try to parse container labels as a map
		// else try to parse container labels as string entries
		containerLabels := map[string]string{}
		if labelMap, ok := spec.ContainerLabels.(map[interface{}]interface{}); ok {
			for k, v := range labelMap {
				containerLabels[k.(string)] = v.(string)
			}
		} else if labelList, ok := spec.ContainerLabels.([]interface{}); ok {
			for _, s := range labelList {
				a := strings.Split(s.(string), "=")
				containerLabels[a[0]] = a[1]
			}
		}
		publishSpecs := []*service.PublishSpec{}
		for _, p := range spec.Public {
			publishSpecs = append(publishSpecs, &service.PublishSpec{
				Name:         p.Name,
				Protocol:     p.Protocol,
				PublishPort:  p.PublishPort,
				InternalPort: p.InternalPort,
			})
		}
		// add service mode and replicas to spec
		var swarmMode service.SwarmMode
		replicas := spec.Replicas
		mode := spec.Mode

		// supply a default value for mode
		if mode == "" {
			mode = "replicated"
		}

		switch mode {
		case "replicated":
			if replicas < 1 {
				// if replicated then must have at least 1 replica
				replicas = 1
			}
			swarmMode = &service.ServiceSpec_Replicated{
				Replicated: &service.ReplicatedService{Replicas: replicas},
			}
		case "global":
			if replicas != 0 {
				// global mode can't specify replicas (only allowed 1 per node)
				return fmt.Errorf("replicas can only be used with replicated mode")
			}
			swarmMode = &service.ServiceSpec_Global{
				Global: &service.GlobalService{},
			}
		default:
			return fmt.Errorf("invalid option for mode: %s", mode)
		}

		// add custom network connection
		networkAttachment := []*service.NetworkAttachment{}
		if spec.Networks != nil {
			for name, data := range spec.Networks {
				trueName, ok := networkMap[name]
				if !ok {
					trueName = name
				}
				networkAttachment = append(networkAttachment, &service.NetworkAttachment{
					Target:  trueName,
					Aliases: data.Aliases,
				})
			}
		}

		stack.Services = append(stack.Services, &service.ServiceSpec{
			Name:            name,
			Image:           spec.Image,
			PublishSpecs:    publishSpecs,
			Mode:            swarmMode,
			Env:             env,
			Args:            args,
			Labels:          labels,
			ContainerLabels: containerLabels,
			Networks:        networkAttachment,
			Mounts:          spec.Mounts,
		})
	}
	return nil
}
