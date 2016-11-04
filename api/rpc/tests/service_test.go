package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/stretchr/testify/assert"
)

var service1 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	Mode: &service.ServiceSpec_Replicated{
		Replicated: &service.ReplicatedService{Replicas: 2},
	},
}

var service2 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	Labels: map[string]string{
		"label1": "value1",
		"label2": "value2",
	},
}

var service3 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	ContainerLabels: map[string]string{
		"label1": "value1",
		"label2": "value2",
	},
}

var service4 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	Env: []string{
		"Var1=value1",
		"Var2=value2",
	},
}

var service5 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	PublishSpecs: []*service.PublishSpec{
		{
			PublishPort:  3001,
			InternalPort: 3000,
		},
	},
}

var service6 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	PublishSpecs: []*service.PublishSpec{
		{
			Name:         "www",
			InternalPort: 80,
		},
	},
}

var service7 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	Args: []string{
		"arg1=value1",
		"arg2=value2",
	},
}

var service8 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	Mounts: []string{
		"/tmp",
	},
}

var service9 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	Mounts: []string{
		"/tmp:/tmp2",
	},
}

var service10 = service.ServiceSpec{
	Image: "appcelerator/pinger",
	Mounts: []string{
		"essai:/tmp2",
	},
}

var serviceList = []*service.ServiceSpec{
	&service1,
	&service2,
	&service3,
	&service4,
	&service5,
	&service6,
	&service7,
	&service8,
	&service9,
	&service10,
}

//Test two stacks life cycle in the same time
func TestServices(t *testing.T) {
	for i, serv := range serviceList {
		name := fmt.Sprintf("service-test%d-%d", i+1, time.Now().Unix())
		serv.Name = name
		respc, errc := serviceClient.Create(ctx, &service.ServiceCreateRequest{
			ServiceSpec: serv,
		})
		if errc != nil {
			t.Fatal(errc)
		}
		assert.NotEmpty(t, respc.Id, "returned service id should not be empty after create")
		respr, errr := serviceClient.Remove(ctx, &service.RemoveRequest{
			Ident: respc.Id,
		})
		if errr != nil {
			t.Fatal(errr)
		}
		assert.NotEmpty(t, respr.Ident, "returned service id should not be empty after remove")
	}
}
