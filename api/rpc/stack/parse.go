package stack

import (
	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v2"
	"strings"
)

type serviceMap struct {
	Image       string      `yaml:"image"`
	Ports       []string    `yaml:"ports"`
	Replicas    int         `yaml:"replicas"`
	Environment interface{} `yaml:"environment"`
	Labels      interface{} `yaml:"labels"`
}

// Service represents a docker service for use in a docker stack
type Service struct {
	Name        string
	Image       string
	Ports       map[nat.Port][]nat.PortBinding
	Replicas    int
	Environment map[string]string
	Labels      map[string]string
}

func parseStackYaml(in string) (out []Service, err error) {
	out = []Service{}
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
		_, ports, err := nat.ParsePortSpecs(d.Ports)
		if err != nil {
			return nil, err
		}
		r := d.Replicas
		if r == 0 {
			r = 1
		}
		out = append(out, Service{
			Name:        n,
			Image:       d.Image,
			Ports:       ports,
			Replicas:    r,
			Environment: e,
			Labels:      l,
		})
	}
	return
}

func parseAsServiceMap(b []byte) (out map[string]serviceMap, err error) {
	err = yaml.Unmarshal(b, &out)
	return
}
