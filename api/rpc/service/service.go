package service

import (
	"strings"

	"github.com/appcelerator/amp/pkg/docker"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server is used to implement log.LogServer
type Server struct {
	Docker *docker.Docker
}

// Tasks implements service.Containers
func (s *Server) Tasks(ctx context.Context, in *TasksRequest) (*TasksReply, error) {
	list, err := s.Docker.GetClient().TaskList(ctx, types.TaskListOptions{})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	taskList := &TasksReply{}
	for _, item := range list {
		if strings.HasPrefix(item.ServiceID, in.ServiceId) || item.Name == in.ServiceId {
			task := &Task{
				Id:           item.ID,
				Image:        strings.Split(item.Spec.ContainerSpec.Image, "@")[0],
				State:        string(item.Status.State),
				DesiredState: string(item.DesiredState),
				NodeId:       item.NodeID,
			}
			taskList.Tasks = append(taskList.Tasks, task)
		}
	}
	return taskList, nil
}
