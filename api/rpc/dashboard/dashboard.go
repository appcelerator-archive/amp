package dashboard

import (
	"log"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/dashboards"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement dashboard.Server
type Server struct {
	Dashboards dashboards.Interface
}

func convertError(err error) error {
	switch err {
	case dashboards.InvalidName:
		return status.Errorf(codes.InvalidArgument, err.Error())
	case dashboards.AlreadyExists:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case dashboards.NotFound:
		return status.Errorf(codes.NotFound, err.Error())
	case accounts.NotAuthorized:
		return status.Errorf(codes.PermissionDenied, err.Error())
	}
	return status.Errorf(codes.Internal, err.Error())
}

// Create implements dashboard.Server
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	dashboard, err := s.Dashboards.Create(ctx, in.Name, in.Data)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully created dashboard:", dashboard)
	return &CreateReply{Dashboard: dashboard}, nil
}

// Get implements account.GetDashboard
func (s *Server) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	// Get the dashboard
	dashboard, err := s.Dashboards.Get(ctx, in.Id)
	if err != nil {
		return nil, convertError(err)
	}
	if dashboard == nil {
		return nil, status.Errorf(codes.NotFound, "dashboard not found: %s", in.Id)
	}
	log.Println("Successfully retrieved dashboard", dashboard.Name)
	return &GetReply{Dashboard: dashboard}, nil
}

// List implements dashboard.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	dashboards, err := s.Dashboards.List(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully listed dashboards")
	return &ListReply{Dashboards: dashboards}, nil
}

// UpdateName implements dashboard.Server
func (s *Server) UpdateName(ctx context.Context, in *UpdateNameRequest) (*empty.Empty, error) {
	if err := s.Dashboards.UpdateName(ctx, in.Id, in.Name); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Printf("Successfully renamed dashboard %s to %s\n", in.Id, in.Name)
	return &empty.Empty{}, nil
}

// UpdateData implements dashboard.Server
func (s *Server) UpdateData(ctx context.Context, in *UpdateDataRequest) (*empty.Empty, error) {
	if err := s.Dashboards.UpdateData(ctx, in.Id, in.Data); err != nil {
		return &empty.Empty{}, convertError(err)
	}
	log.Printf("Successfully updated dashboard %d data\n", in.Id)
	return &empty.Empty{}, nil
}

// Remove implements dashboard.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*empty.Empty, error) {
	if err := s.Dashboards.Delete(ctx, in.Id); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully deleted dashboard", in.Id)
	return &empty.Empty{}, nil
}
