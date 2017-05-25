package resource

import (
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/stacks"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement resource.ResourceServer
type Server struct {
	Stacks stacks.Interface
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

// ListResources implements resource.ListResources
func (s *Server) ListResources(ctx context.Context, in *ListResourcesRequest) (*ListResourcesReply, error) {
	stacks, err := s.Stacks.ListStacks(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	reply := &ListResourcesReply{}
	user := auth.GetUser(ctx)
	for _, stack := range stacks {
		if stack.Owner.Name == user {
			reply.Resources = append(reply.Resources, &ResourceEntry{Id: stack.Id, Type: ResourceType_RESOURCE_STACK, Name: stack.Name})
		}
	}
	return reply, nil
}
