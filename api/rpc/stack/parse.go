package stack

import (
	"strings"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/docker/docker/pkg/stringid"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

type serviceSpec struct {
	Image       string        `yaml:"image"`
	Replicas    uint64        `yaml:"replicas"`
	Environment interface{}   `yaml:"environment"`
	Labels      interface{}   `yaml:"labels"`
	Public      []publishSpec `yaml:"public"`
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

		replicas := spec.Replicas
		if replicas == 0 {
			replicas = 1
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

		stack.Services = append(stack.Services, &service.ServiceSpec{
			Name:         name,
			Image:        spec.Image,
			Replicas:     replicas,
			Env:          env,
			Labels:       labels,
			PublishSpecs: publishSpecs,
		})
	}

	return
}

func parseServiceMap(b []byte) (out map[string]serviceSpec, err error) {
	err = yaml.Unmarshal(b, &out)
	return
}
