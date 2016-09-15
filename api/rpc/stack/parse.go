package stack

import (
	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v2"
	"strings"
)

type serviceMap struct {
	Image       string      `yaml:"image"`
	Ports       []string    `yaml:"ports"`
	Replicas    uint64      `yaml:"replicas"`
	Environment interface{} `yaml:"environment"`
	Labels      interface{} `yaml:"labels"`
}

func parseStackYaml(in string) (out *Stack, err error) {
	out = &Stack{}
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
		_, natPorts, err := nat.ParsePortSpecs(d.Ports)
		if err != nil {
			return &Stack{}, err
		}
		ports := map[string]*Service_PortBindings{}
		for p, bs := range natPorts {
			pbs := Service_PortBindings{}
			for _, b := range bs {
				pbs.PortBindings = append(pbs.PortBindings, &Service_PortBinding{
					HostIp:   b.HostIP,
					HostPort: b.HostPort,
				})
			}
			ports[string(p)] = &pbs
		}
		r := d.Replicas
		if r == 0 {
			r = 1
		}
		out.Services = append(out.Services, &Service{
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
