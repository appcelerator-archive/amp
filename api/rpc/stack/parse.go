package stack

import (
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/docker/docker/pkg/stringid"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

type serviceSpec struct {
	Image           string        `yaml:"image"`
	Public          []publishSpec `yaml:"public"`
	Mode            string        `yaml:"mode"`
	Replicas        uint64        `yaml:"replicas"`
	Environment     interface{}   `yaml:"environment"`
	Labels          interface{}   `yaml:"labels"`
	ContainerLabels interface{}   `yaml:"container_labels"`
}

type publishSpec struct {
	Name         string `yaml:"name"`
	Protocol     string `yaml:"protocol"`
	PublishPort  uint32 `yaml:"publish_port"`
	InternalPort uint32 `yaml:"internal_port"`
}

// ParseStackfile create a new stack from yaml
func ParseStackfile(ctx context.Context, in string) (stack *Stack, err error) {
	stack = &Stack{}
	stack.Id = stringid.GenerateNonCryptoID()
	serviceMap, err := parseServiceMap([]byte(in))
	if err != nil {
		return
	}
	for name, spec := range serviceMap {
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
		labels := map[string]string{}
		if labelMap, ok := spec.Labels.(map[interface{}]interface{}); ok {
			for k, v := range labelMap {
				labels[k.(string)] = v.(string)
			}
		} else if labelList, ok := spec.Labels.([]interface{}); ok {
			for _, s := range labelList {
				a := strings.Split(s.(string), "=")
				k := a[0]
				v := a[1]
				labels[k] = v
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
				k := a[0]
				v := a[1]
				containerLabels[k] = v
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
				err = fmt.Errorf("replicas can only be used with replicated mode")
				return
			}
			swarmMode = &service.ServiceSpec_Global{
				Global: &service.GlobalService{},
			}
		default:
			err = fmt.Errorf("invalid option for mode: %s", mode)
			return
		}

		stack.Services = append(stack.Services, &service.ServiceSpec{
			Name:            name,
			Image:           spec.Image,
			PublishSpecs:    publishSpecs,
			Mode:            swarmMode,
			Env:             env,
			Labels:          labels,
			ContainerLabels: containerLabels,
		})
	}
	return
}

func parseServiceMap(b []byte) (out map[string]serviceSpec, err error) {
	err = yaml.Unmarshal(b, &out)
	return
}
