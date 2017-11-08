package docker

import (
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
	"golang.org/x/net/context"
)

type TasksBySlot []swarm.Task

func (t TasksBySlot) Len() int {
	return len(t)
}

func (t TasksBySlot) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TasksBySlot) Less(i, j int) bool {
	// Sort by slot.
	if t[i].Slot != t[j].Slot {
		return t[i].Slot < t[j].Slot
	}

	// If same slot, sort by most recent.
	return t[j].Meta.CreatedAt.Before(t[i].CreatedAt)
}

// TaskList list the tasks
func (d *Docker) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) {
	return d.client.TaskList(ctx, options)
}
