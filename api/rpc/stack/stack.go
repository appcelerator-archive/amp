package stack

import (
	"log"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/appcelerator/amp/pkg/docker"
	"golang.org/x/net/context"
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
	case stacks.StackAlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case stacks.StackNotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case accounts.NotAuthorized:
		return status.Errorf(codes.PermissionDenied, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}

// Deploy implements stack.Server
func (s *Server) Deploy(ctx context.Context, in *DeployRequest) (*DeployReply, error) {
	log.Println("[stack] Deploy", in.String())

	// Check if stack already exists
	stack, err := s.Stacks.GetStackByName(ctx, in.Name)
	if err != nil {
		return nil, convertError(err)
	}
	if stack == nil {
		if stack, err = s.Stacks.CreateStack(ctx, in.Name); err != nil {
			return nil, convertError(err)
		}
	}

	// Deploy stack
	output, err := s.Docker.StackDeploy(ctx, stack.Name, in.Compose)
	if err != nil {
		s.Stacks.DeleteStack(ctx, stack.Id)
		return nil, convertError(err)
	}
	return &DeployReply{FullName: stack.Name, Answer: output}, nil
}

// List implements stack.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Println("[stack] List", in.String())

	// List stacks
	output, err := s.Docker.StackList(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	reply := &ListReply{}
	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] {
		if len(strings.Fields(line)) == 0 {
			continue
		}
		entry := s.toStackListEntry(ctx, line)
		if entry == nil {
			continue
		}
		reply.Entries = append(reply.Entries, entry)
	}
	return reply, nil
}

func (s *Server) toStackListEntry(ctx context.Context, line string) *StackListEntry {
	cols := strings.Fields(line)
	name := cols[0]
	services := cols[1]
	stk, err := s.Stacks.GetStackByName(ctx, name)
	if err != nil || stk == nil {
		return nil
	}
	return &StackListEntry{Stack: stk, Services: services}
}

// Remove implements stack.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	log.Println("[stack] Remove", in.String())

	// Retrieve the stack
	stack, dockerErr := s.Stacks.GetStackByFragmentOrName(ctx, in.Stack)
	if dockerErr != nil {
		return nil, convertError(dockerErr)
	}

	// Check authorization
	if !s.Accounts.IsAuthorized(ctx, stack.Owner, accounts.DeleteAction, accounts.StackRN, stack.Id) {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	// Remove stack
	output, dockerErr := s.Docker.StackRemove(ctx, stack.Name)
	storageErr := s.Stacks.DeleteStack(ctx, stack.Id)
	if dockerErr != nil {
		return nil, convertError(dockerErr)
	}
	if storageErr != nil {
		return nil, convertError(dockerErr)
	}

	log.Printf("Stack %s removed", in.Stack)
	return &RemoveReply{Answer: output}, nil
}
