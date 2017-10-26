package stack

import (
	"strings"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/appcelerator/amp/pkg/docker"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement stack.Server
type Server struct {
	Accounts accounts.Interface
	Docker   *docker.Docker
	Stacks   stacks.Interface
}

func convertError(err error) error {
	switch err {
	case stacks.InvalidName:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case stacks.AlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case stacks.NotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case accounts.NotAuthorized:
		return status.Errorf(codes.PermissionDenied, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}

// Deploy implements stack.Server
func (s *Server) Deploy(ctx context.Context, in *DeployRequest) (*DeployReply, error) {
	// Check if stack is using restricted resources
	compose, err := s.Docker.ComposeParse(ctx, in.Compose, in.Environment)
	if err != nil {
		return nil, convertError(err)
	}
	if !s.Docker.ComposeIsAuthorized(compose) {
		return nil, status.Errorf(codes.FailedPrecondition, "This compose file requires access to reserved AMP resources.")
	}

	// Check if stack already exists
	stack, err := s.Stacks.GetByName(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}

	if stack == nil {
		if stack, err = s.Stacks.Create(ctx, in.Name); err != nil {
			return nil, convertError(err)
		}
	} else {
		// if it does, make sure we have the right to update it
		if !s.Accounts.IsAuthorized(ctx, stack.Owner, accounts.UpdateAction, accounts.StackRN, stack.Id) {
			return nil, stacks.AlreadyExists
		}
	}

	// Deploy stack
	output, err := s.Docker.StackDeploy(ctx, stack.Name, in.Compose, in.Config, in.Environment)
	if err != nil {
		s.Stacks.Delete(ctx, stack.Id)
		return nil, convertError(err)
	}
	log.Infoln("Successfully deployed stack:", stack)
	return &DeployReply{Id: stack.Id, FullName: stack.Name, Answer: output}, nil
}

// List implements stack.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Infoln("[stack] List", in.String())

	// List stacks
	reply := &ListReply{}
	stacks, err := s.Stacks.List(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	for _, stack := range stacks {
		entry, err := s.toStackListEntry(ctx, stack)
		if err != nil {
			return nil, convertError(err)
		}
		reply.Entries = append(reply.Entries, entry)
	}
	return reply, nil
}

func (s *Server) toStackListEntry(ctx context.Context, stack *stacks.Stack) (*StackListEntry, error) {
	status, err := s.Docker.StackStatus(ctx, stack.Name)
	if err != nil {
		return nil, convertError(err)
	}
	log.Infoln("[stack] Stack", stack.Name, "is", status.Status, "with", status.RunningServices, "out of", status.TotalServices, "services and", status.FailedServices, "failed services")
	return &StackListEntry{Stack: stack, RunningServices: status.RunningServices, FailedServices: status.FailedServices, TotalServices: status.TotalServices, Status: status.Status}, nil
}

// Remove implements stack.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	log.Infoln("[stack] Remove", in.String())

	// Retrieve the stack
	stack, dockerErr := s.Stacks.GetByFragmentOrName(ctx, in.Stack)
	if dockerErr != nil {
		return nil, convertError(dockerErr)
	}
	if stack == nil {
		return nil, stacks.NotFound
	}

	// Check authorization
	if !s.Accounts.IsAuthorized(ctx, stack.Owner, accounts.DeleteAction, accounts.StackRN, stack.Id) {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	// Remove stack
	output, dockerErr := s.Docker.StackRemove(ctx, stack.Name)
	storageErr := s.Stacks.Delete(ctx, stack.Id)
	if dockerErr != nil {
		return nil, convertError(dockerErr)
	}
	if storageErr != nil {
		return nil, convertError(storageErr)
	}

	log.Infof("Stack %s removed", in.Stack)
	return &RemoveReply{Answer: output}, nil
}

// Services Ctx implements stack.Server
func (s *Server) Services(ctx context.Context, in *ServicesRequest) (*ServicesReply, error) {
	log.Infoln("[stack] Services", in.String())

	stack, err := s.Stacks.GetByFragmentOrName(ctx, in.StackName)
	if err != nil {
		return nil, convertError(err)
	}
	if stack == nil {
		return nil, stacks.NotFound
	}

	output, dockerErr := s.Docker.StackServices(ctx, stack.Name, false)
	if dockerErr != nil {
		log.Infof("error : %v\n", dockerErr)
		return nil, convertError(dockerErr)
	}

	cols := strings.Split(output, "\n")
	ans := &ServicesReply{
		Services: []*StackService{},
	}
	for _, col := range cols[1:] {
		service := s.getOneServiceListLine(ctx, col)
		if service != nil {
			ans.Services = append(ans.Services, service)
		}
	}
	return ans, nil
}

func (s *Server) getOneServiceListLine(ctx context.Context, line string) *StackService {
	cols := strings.Split(line, " ")
	nn := 0
	service := &StackService{}
	for _, val := range cols {
		if val != "" {
			nn++
			switch nn {
			case 1:
				service.Id = val
				break
			case 2:
				service.Name = val
				break
			case 3:
				service.Mode = val
				break
			case 4:
				service.Replicas = val
				break
			case 5:
				service.Image = val
				break
			}
		}
	}
	return service
}
