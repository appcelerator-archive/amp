package service

import (
	"fmt"
	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/docker/cli/cli/command/formatter"
)

// GetServicesStatus returns a map of mode and replicas
func GetServicesStatus(services []swarm.Service, nodes []swarm.Node, tasks []swarm.Task) map[string]formatter.ServiceListInfo {
	running := map[string]int{}
	tasksNoShutdown := map[string]int{}

	activeNodes := make(map[string]struct{})
	for _, n := range nodes {
		if n.Status.State != swarm.NodeStateDown {
			activeNodes[n.ID] = struct{}{}
		}
	}

	for _, task := range tasks {
		if task.DesiredState != swarm.TaskStateShutdown {
			tasksNoShutdown[task.ServiceID]++
		}

		if _, nodeActive := activeNodes[task.NodeID]; nodeActive && task.Status.State == swarm.TaskStateRunning {
			running[task.ServiceID]++
		}
	}

	info := map[string]formatter.ServiceListInfo{}
	for _, service := range services {
		info[service.ID] = formatter.ServiceListInfo{}
		if service.Spec.Mode.Replicated != nil && service.Spec.Mode.Replicated.Replicas != nil {
			info[service.ID] = formatter.ServiceListInfo{
				Mode:     "replicated",
				Replicas: fmt.Sprintf("%d/%d", running[service.ID], *service.Spec.Mode.Replicated.Replicas),
			}
		} else if service.Spec.Mode.Global != nil {
			info[service.ID] = formatter.ServiceListInfo{
				Mode:     "global",
				Replicas: fmt.Sprintf("%d/%d", running[service.ID], tasksNoShutdown[service.ID]),
			}
		}
	}
	return info
}
