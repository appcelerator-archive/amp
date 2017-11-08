package service

import (
	"encoding/json"
	"sort"
	"strings"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement ServiceServer
type Server struct {
	Accounts accounts.Interface
	Docker   *docker.Docker
	Stacks   stacks.Interface
}

// Service constants
const (
	RoleLabel          = "io.amp.role"
	LatestTag          = "latest"
	GlobalMode         = "global"
	ReplicatedMode     = "replicated"
	StackNameLabelName = "com.docker.stack.namespace"
)

// Ps implements service.Ps
func (s *Server) Ps(ctx context.Context, in *PsRequest) (*PsReply, error) {
	log.Infoln("[service] PsService", in.Service)
	args := filters.NewArgs()
	args.Add("service", in.Service)
	tasks, err := s.Docker.TaskList(ctx, types.TaskListOptions{Filters: args})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	// Sort tasks by slot, then by most recent
	sort.Stable(docker.TasksBySlot(tasks))
	taskList := &PsReply{}

	for _, task := range tasks {
		task := &Task{
			Id:           task.ID,
			Image:        strings.Split(task.Spec.ContainerSpec.Image, "@")[0],
			CurrentState: strings.ToUpper(string(task.Status.State)),
			DesiredState: strings.ToUpper(string(task.DesiredState)),
			NodeId:       task.NodeID,
			Error:        task.Status.Err,
			Slot:         int32(task.Slot),
		}
		taskList.Tasks = append(taskList.Tasks, task)
	}
	return taskList, nil
}

// List implements service.List
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Infoln("[service] List ", in.Stack)
	serviceList, err := s.Docker.ServicesList(ctx, types.ServiceListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	reply := &ListReply{}
	for _, service := range serviceList {
		if _, ok := service.Spec.Labels[RoleLabel]; ok {
			continue // ignore amp infrastructure services
		}
		stackName := service.Spec.Labels[StackNameLabelName]
		if in.Stack != "" && stackName != in.Stack {
			continue // filter based on provided stack name
		}
		entry := &ServiceEntry{
			Id:   service.ID,
			Name: service.Spec.Name,
		}
		image := service.Spec.TaskTemplate.ContainerSpec.Image
		if strings.Contains(image, "@") {
			image = strings.Split(image, "@")[0] // trimming the hash
		}
		entry.Image = image
		entry.Tag = LatestTag
		if strings.Contains(image, ":") {
			index := strings.LastIndex(image, ":")
			entry.Image = image[:index]
			entry.Tag = image[index+1:]
		}
		entry.Mode = ReplicatedMode
		if service.Spec.Mode.Global != nil {
			entry.Mode = GlobalMode
		}
		status, err := s.Docker.ServiceStatus(ctx, &service)
		if err != nil {
			return nil, err
		}
		entry.RunningTasks = status.RunningTasks
		entry.TotalTasks = status.TotalTasks
		entry.Status = status.Status

		reply.Entries = append(reply.Entries, entry)
	}
	return reply, nil
}

// Inspect inspects a service
func (s *Server) Inspect(ctx context.Context, in *InspectRequest) (*InspectReply, error) {
	log.Infoln("[service] Inspect", in.Service)
	serviceEntity, err := s.Docker.ServiceInspect(ctx, in.Service)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	entity, _ := json.MarshalIndent(serviceEntity, "", "	")
	return &InspectReply{Json: string(entity)}, nil
}

// Scale scales a service
func (s *Server) Scale(ctx context.Context, in *ScaleRequest) (*empty.Empty, error) {
	log.Infoln("[service] Scale", in.Service)
	serviceEntity, err := s.Docker.ServiceInspect(ctx, in.Service)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	stackName := serviceEntity.Spec.Labels[StackNameLabelName]

	stack, dockerErr := s.Stacks.GetByFragmentOrName(ctx, stackName)
	if dockerErr != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	if stack == nil {
		return nil, stacks.NotFound
	}

	// Check authorization
	if !s.Accounts.IsAuthorized(ctx, stack.Owner, accounts.UpdateAction, accounts.StackRN, stack.Id) {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	if err := s.Docker.ServiceScale(ctx, in.Service, in.Replicas); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &empty.Empty{}, nil
}
