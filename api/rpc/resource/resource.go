package resource

import (
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/dashboards"
	"github.com/appcelerator/amp/data/stacks"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement resource.ResourceServer
type Server struct {
	Stacks     stacks.Interface
	Dashboards dashboards.Interface
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

// ListResources implements resource.ListResources
func (s *Server) ListResources(ctx context.Context, in *ListResourcesRequest) (*ListResourcesReply, error) {
	reply := &ListResourcesReply{}

	account := auth.GetActiveOrganization(ctx)
	if account == "" {
		account = auth.GetUser(ctx)
	}

	// Stacks
	stacks, err := s.Stacks.ListStacks(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	for _, stack := range stacks {
		if stack.Owner.Name == account {
			reply.Resources = append(reply.Resources, &ResourceEntry{Id: stack.Id, Type: ResourceType_RESOURCE_STACK, Name: stack.Name})
		}
	}

	// Dashboards
	dashboards, err := s.Dashboards.List(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	for _, dashboard := range dashboards {
		if dashboard.Owner.Name == account {
			reply.Resources = append(reply.Resources, &ResourceEntry{Id: dashboard.Id, Type: ResourceType_RESOURCE_DASHBOARD, Name: dashboard.Name})
		}
	}
	return reply, nil
}
