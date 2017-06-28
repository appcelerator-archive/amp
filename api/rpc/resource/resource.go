package resource

import (
	log "github.com/Sirupsen/logrus"
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/dashboards"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement resource.ResourceServer
type Server struct {
	Accounts   accounts.Interface
	Dashboards dashboards.Interface
	Stacks     stacks.Interface
}

func convertError(err error) error {
	switch err {
	case stacks.InvalidName, dashboards.InvalidName:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case stacks.AlreadyExists, dashboards.AlreadyExists, accounts.ResourceAlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case stacks.NotFound, dashboards.NotFound, accounts.ResourceNotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case accounts.NotAuthorized:
		return status.Errorf(codes.PermissionDenied, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}

// List implements resource.List
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	reply := &ListReply{}

	activeOrganization := auth.GetActiveOrganization(ctx)
	if activeOrganization == "" {
		return reply, nil
	}

	// Stacks
	stacks, err := s.Stacks.List(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	for _, stack := range stacks {
		if stack.Owner.Organization == activeOrganization {
			reply.Resources = append(reply.Resources, &ResourceEntry{Id: stack.Id, Type: ResourceType_RESOURCE_STACK, Name: stack.Name, Owner: stack.Owner})
		}
	}

	// Dashboards
	dashboards, err := s.Dashboards.List(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	for _, dashboard := range dashboards {
		if dashboard.Owner.Organization == activeOrganization {
			reply.Resources = append(reply.Resources, &ResourceEntry{Id: dashboard.Id, Type: ResourceType_RESOURCE_DASHBOARD, Name: dashboard.Name, Owner: dashboard.Owner})
		}
	}
	log.Infoln("Successfully listed resources for organization", activeOrganization)
	return reply, nil
}

// AddToTeam implements resource.AddToTeam
func (s *Server) AddToTeam(ctx context.Context, in *AddToTeamRequest) (*empty.Empty, error) {
	reply, err := s.List(ctx, &ListRequest{})
	if err != nil {
		return &empty.Empty{}, convertError(err)
	}

	// Make sure the resource belongs to the given organization
	var resource *ResourceEntry
	for _, res := range reply.Resources {
		if res.Id == in.ResourceId {
			resource = res
			break
		}
	}
	if resource == nil {
		return &empty.Empty{}, status.Errorf(codes.NotFound, "Resource doesn't belong to the given organization")
	}

	RN := ""
	switch resource.Type {
	case ResourceType_RESOURCE_DASHBOARD:
		RN = accounts.DashboardRN
	case ResourceType_RESOURCE_STACK:
		RN = accounts.StackRN
	default:
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "Resource type is not supported")
	}

	// Check authorization over resource
	if !s.Accounts.IsAuthorized(ctx, resource.Owner, accounts.AdminAction, RN, in.ResourceId) {
		return &empty.Empty{}, convertError(accounts.NotAuthorized)
	}

	// Add resource to the team
	if err := s.Accounts.AddResourceToTeam(ctx, in.OrganizationName, in.TeamName, in.ResourceId); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Infof("Successfully added resource %s to team %s in organization %s\n", in.ResourceId, in.TeamName, in.OrganizationName)
	return &empty.Empty{}, nil
}

// ChangePermissionLevel implements resource.ChangePermissionLevel
func (s *Server) ChangePermissionLevel(ctx context.Context, in *ChangePermissionLevelRequest) (*empty.Empty, error) {
	reply, err := s.List(ctx, &ListRequest{})
	if err != nil {
		return &empty.Empty{}, convertError(err)
	}

	// Make sure the resource belongs to the given organization
	var resource *ResourceEntry
	for _, res := range reply.Resources {
		if res.Id == in.ResourceId {
			resource = res
			break
		}
	}
	if resource == nil {
		return &empty.Empty{}, status.Errorf(codes.NotFound, "Resource doesn't belong to the given organization")
	}

	RN := ""
	switch resource.Type {
	case ResourceType_RESOURCE_DASHBOARD:
		RN = accounts.DashboardRN
	case ResourceType_RESOURCE_STACK:
		RN = accounts.StackRN
	default:
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "Resource type is not supported")
	}

	// Check authorization over resource
	if !s.Accounts.IsAuthorized(ctx, resource.Owner, accounts.AdminAction, RN, in.ResourceId) {
		return &empty.Empty{}, convertError(accounts.NotAuthorized)
	}

	if err := s.Accounts.ChangeTeamResourcePermissionLevel(ctx, in.OrganizationName, in.TeamName, in.ResourceId, in.PermissionLevel); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Infof("Successfully changed permission level over resource %s in team %s to %s\n", in.ResourceId, in.TeamName, in.PermissionLevel.String())
	return &empty.Empty{}, nil
}

// RemoveFromTeam implements resource.RemoveFromTeam
func (s *Server) RemoveFromTeam(ctx context.Context, in *RemoveFromTeamRequest) (*empty.Empty, error) {
	if err := s.Accounts.RemoveResourceFromTeam(ctx, in.OrganizationName, in.TeamName, in.ResourceId); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Infof("Successfully removed resource %s from teams %s in organization %s\n", in.ResourceId, in.TeamName, in.OrganizationName)
	return &empty.Empty{}, nil
}
