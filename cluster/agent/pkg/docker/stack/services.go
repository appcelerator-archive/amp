package stack

import (
	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type ServicesOptions struct {
	quiet     bool
	format    string
	filter    opts.FilterOpt
	namespace string
}

func ListServices(ctx context.Context, client client.APIClient, options types.ServiceListOptions) ([]swarm.Service, error) {
	return client.ServiceList(ctx, options)
}
