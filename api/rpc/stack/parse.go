package stack

import (
	"strings"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/docker/docker/pkg/stringid"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"path"
)

type serviceMap struct {
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

// NewStackfromYaml create a new stack from yaml
func NewStackFromYaml(ctx context.Context, in string) (stack *Stack, err error) {
	stack = &Stack{}
	stack.Id = stringid.GenerateNonCryptoID()
	b := []byte(in)
	sm, err := parseAsServiceMap(b)
	if err != nil {
		return
	}
	for n, d := range sm {
		e := map[string]string{}
		l := map[string]string{}
		em, ok := d.Environment.(map[interface{}]interface{})
		if ok {
			for k, v := range em {
				e[k.(string)] = v.(string)
			}
		}
		ea, ok := d.Environment.([]interface{})
		if ok {
			for _, s := range ea {
				a := strings.Split(s.(string), "=")
				k := a[0]
				v := a[1]
				e[k] = v
			}
		}
		lm, ok := d.Labels.(map[interface{}]interface{})
		if ok {
			for k, v := range lm {
				l[k.(string)] = v.(string)
			}
		}
		la, ok := d.Labels.([]interface{})
		if ok {
			for _, s := range la {
				a := strings.Split(s.(string), "=")
				k := a[0]
				v := a[1]
				l[k] = v
			}
		}
		r := d.Replicas
		if r == 0 {
			r = 1
		}
		publishSpecs := []*service.PublishSpec{}
		for _, p := range d.Public {
			publishSpecs = append(publishSpecs, &service.PublishSpec{
				Name:         p.Name,
				Protocol:     p.Protocol,
				PublishPort:  p.PublishPort,
				InternalPort: p.InternalPort,
			})
		}
		stack.Services = append(stack.Services, &service.ServiceSpec{
			Name:         n,
			Image:        d.Image,
			Replicas:     r,
			Env:          e,
			Labels:       l,
			PublishSpecs: publishSpecs,
		})
	}

	// Store stack
	if err = runtime.Store.Create(ctx, path.Join("stacks", stack.Id), stack, nil, 0); err != nil {
		return
	}

	// Create stack state
	if err = stackStateMachine.CreateState(stack.Id, int32(StackState_Stopped)); err != nil {
		return
	}
	return
}

func parseAsServiceMap(b []byte) (out map[string]serviceMap, err error) {
	err = yaml.Unmarshal(b, &out)
	return
}
