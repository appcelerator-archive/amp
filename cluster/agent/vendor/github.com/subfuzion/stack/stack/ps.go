package stack

import (
	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type PsOptions struct {
	filter    opts.FilterOpt
	noTrunc   bool
	namespace string
	noResolve bool
	quiet     bool
	format    string
}

func ListTasks(ctx context.Context, client client.APIClient, options types.TaskListOptions) ([]swarm.Task, error) {
	return client.TaskList(ctx, options)
}
