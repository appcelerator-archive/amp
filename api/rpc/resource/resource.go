package resource

import (
	"strings"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/dashboards"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/elastic/go-lumber/log"
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
	case stacks.AlreadyExists, dashboards.AlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case stacks.NotFound, dashboards.NotFound:
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
	stacks, err := s.Stacks.ListStacks(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	for _, stack := range stacks {
		if stack.Owner.Name == activeOrganization {
			reply.Resources = append(reply.Resources, &ResourceEntry{Id: stack.Id, Type: ResourceType_RESOURCE_STACK, Name: stack.Name})
		}
	}

	// Dashboards
	dashboards, err := s.Dashboards.List(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	for _, dashboard := range dashboards {
		if dashboard.Owner.Name == activeOrganization {
			reply.Resources = append(reply.Resources, &ResourceEntry{Id: dashboard.Id, Type: ResourceType_RESOURCE_DASHBOARD, Name: dashboard.Name})
		}
	}
	log.Println("Successfully listed resources for organization", activeOrganization)
	return reply, nil
}

func (s *Server) isAuthorized(ctx context.Context, request *IsAuthorizedRequest) bool {
	var owner *accounts.Account
	var action, resourceType string
	resourceID := request.Id
	switch request.Type {
	case ResourceType_RESOURCE_USER:
		resourceType = accounts.UserRN
		owner = &accounts.Account{
			Type: accounts.AccountType_USER,
			Name: request.Id,
		}
	case ResourceType_RESOURCE_ORGANIZATION:
		resourceType = accounts.OrganizationRN
		owner = &accounts.Account{
			Type: accounts.AccountType_ORGANIZATION,
			Name: request.Id,
		}
	case ResourceType_RESOURCE_TEAM:
		resourceType = accounts.TeamRN
		IDs := strings.Split(request.Id, "/") // the team ID is the concatenation of orgName/teamName
		switch len(IDs) {
		case 1:
			owner = &accounts.Account{
				Type: accounts.AccountType_ORGANIZATION,
				Name: IDs[0],
			}
			resourceID = IDs[0]
		case 2:
			owner = &accounts.Account{
				Type: accounts.AccountType_ORGANIZATION,
				Name: IDs[0],
			}
			resourceID = IDs[1]
		default:
			return false
		}
	case ResourceType_RESOURCE_DASHBOARD:
		resourceType = accounts.DashboardRN
		dashboard, err := s.Dashboards.Get(ctx, request.Id)
		if err != nil {
			return false
		}
		if dashboard == nil {
			return false
		}
		owner = dashboard.Owner
	case ResourceType_RESOURCE_STACK:
		resourceType = accounts.StackRN
		stack, err := s.Stacks.GetStackByName(ctx, request.Id)
		if err != nil {
			return false
		}
		if stack == nil {
			return false
		}
		owner = stack.Owner
	}

	switch request.Action {
	case Action_ACTION_CREATE:
		action = accounts.CreateAction
	case Action_ACTION_READ:
		action = accounts.ReadAction
	case Action_ACTION_UPDATE:
		action = accounts.UpdateAction
	case Action_ACTION_DELETE:
		action = accounts.DeleteAction
	}

	return s.Accounts.IsAuthorized(ctx, owner, action, resourceType, resourceID)
}

// Authorizations implements resource.Authorizations
func (s *Server) Authorizations(ctx context.Context, in *AuthorizationsRequest) (*AuthorizationsReply, error) {
	reply := &AuthorizationsReply{}
	for _, request := range in.Requests {
		reply.Replies = append(reply.Replies, &IsAuthorizedReply{
			Id:         request.Id,
			Type:       request.Type,
			Action:     request.Action,
			Authorized: s.isAuthorized(ctx, request),
		})
	}
	log.Println("Successfully retrieved authorizations")
	return reply, nil
}

// AddToTeam implements resource.AddToTeam
func (s *Server) AddToTeam(ctx context.Context, in *AddToTeamRequest) (*empty.Empty, error) {
	reply, err := s.List(ctx, &ListRequest{})
	if err != nil {
		return &empty.Empty{}, err
	}
	found := false
	for _, res := range reply.Resources {
		if res.Id == in.ResourceId {
			found = true
			break
		}
	}
	if !found {
		return &empty.Empty{}, status.Errorf(codes.NotFound, "Resource not found in the given organization")
	}
	if err := s.Accounts.AddResourceToTeam(ctx, in.OrganizationName, in.TeamName, in.ResourceId); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Printf("Successfully added resource %s to team %s in organization %s\n", in.ResourceId, in.TeamName, in.OrganizationName)
	return &empty.Empty{}, nil
}

// RemoveFromTeam implements resource.RemoveFromTeam
func (s *Server) RemoveFromTeam(ctx context.Context, in *RemoveFromTeamRequest) (*empty.Empty, error) {
	if err := s.Accounts.RemoveResourceFromTeam(ctx, in.OrganizationName, in.TeamName, in.ResourceId); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Printf("Successfully removed resource %s from teams %s in organization %s\n", in.ResourceId, in.TeamName, in.OrganizationName)
	return &empty.Empty{}, nil
}
