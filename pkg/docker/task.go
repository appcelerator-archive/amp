package docker

import (
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"docker.io/go-docker/api/types/swarm"
	"golang.org/x/net/context"
)

// checkTasks returns the running and failing tasks of a service
func (d *Docker) checkTasks(ctx context.Context, service string) (map[string]int, error) {
	args := filters.NewArgs()
	args.Add("service", service)
	taskMap := map[string]int{}
	serviceTasks, err := d.TaskList(ctx, types.TaskListOptions{Filters: args})
	if err != nil {
		return nil, err
	}
	for _, serviceTask := range serviceTasks {
		if serviceTask.Status.State == swarm.TaskStateRejected || serviceTask.Status.State == swarm.TaskStateFailed {
			taskMap[StateError]++
		}
		if serviceTask.Status.State == swarm.TaskStateRunning {
			taskMap[StateRunning]++
		}
	}
	return taskMap, nil
}
