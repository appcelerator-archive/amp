package service

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/appcelerator/amp/pkg/docker"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server is used to implement ServiceServer
type Server struct {
	Docker *docker.Docker
}

// Service constants
const (
	RoleLabel      = "io.amp.role"
	LatestTag      = "latest"
	GlobalMode     = "global"
	ReplicatedMode = "replicated"
)

// Tasks implements service.Containers
func (s *Server) Tasks(ctx context.Context, in *TasksRequest) (*TasksReply, error) {
	log.Println("[service] Tasks", in.ServiceId)
	list, err := s.Docker.TaskList(ctx, types.TaskListOptions{})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	taskList := &TasksReply{}
	for _, item := range list {
		if strings.HasPrefix(item.ServiceID, in.ServiceId) {
			task := &Task{
				Id:           item.ID,
				Image:        strings.Split(item.Spec.ContainerSpec.Image, "@")[0],
				CurrentState: strings.ToUpper(string(item.Status.State)),
				DesiredState: strings.ToUpper(string(item.DesiredState)),
				NodeId:       item.NodeID,
			}
			taskList.Tasks = append(taskList.Tasks, task)
		}
	}
	return taskList, nil
}

// ListService implements service.ListService
func (s *Server) ListService(ctx context.Context, in *ServiceListRequest) (*ServiceListReply, error) {
	log.Println("[service] List")
	serviceList, err := s.Docker.ServicesList(ctx, types.ServiceListOptions{})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	reply := &ServiceListReply{}
	for _, service := range serviceList {
		if _, ok := service.Spec.Labels[RoleLabel]; !ok {
			entity := &ServiceEntity{
				Id:   service.ID,
				Name: service.Spec.Name,
			}
			if strings.Contains(service.Spec.TaskTemplate.ContainerSpec.Image, "@") {
				imageTag := strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, "@")[0]
				it := strings.Split(imageTag, ":")
				entity.Image = it[0]
				entity.Tag = it[1]
			} else if strings.Contains(service.Spec.TaskTemplate.ContainerSpec.Image, ":") {
				it := strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, ":")
				entity.Image = it[0]
				entity.Tag = it[1]
			} else {
				entity.Image = service.Spec.TaskTemplate.ContainerSpec.Image
				entity.Tag = LatestTag
			}
			if service.Spec.Mode.Global != nil {
				entity.Mode = GlobalMode
			} else {
				entity.Mode = ReplicatedMode
			}
			response, err := s.serviceStatusReplicas(ctx, entity)
			if err != nil {
				return nil, grpc.Errorf(codes.Internal, "%v", err)
			}
			reply.Entries = append(reply.Entries, response)

		}
	}
	return reply, nil
}

func (s *Server) serviceStatusReplicas(ctx context.Context, service *ServiceEntity) (*ServiceListEntry, error) {
	statusReplicas, err := s.Docker.ServiceState(ctx, service.Id)
	if err != nil {
		return nil, err
	}
	return &ServiceListEntry{Service: service, ReadyTasks: statusReplicas.RunningTasks, TotalTasks: statusReplicas.TotalTasks, Status: statusReplicas.Status}, nil
}

func (s *Server) InspectService(ctx context.Context, in *ServiceInspectRequest) (*ServiceInspectReply, error) {
	log.Println("[service] Inspect", in.ServiceId)
	serviceEntity, err := s.Docker.ServiceInspect(ctx, in.ServiceId)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	entity, _ := json.MarshalIndent(serviceEntity, "", "	")
	return &ServiceInspectReply{ServiceEntity: string(entity)}, nil
}
