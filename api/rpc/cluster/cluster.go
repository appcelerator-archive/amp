package cluster

import (
	"log"

	"golang.org/x/net/context"
)

// Server is used to implement cluster.ClusterServer
type Server struct {
}

// Create implements cluster.Server
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	log.Println("[cluster] Create", in.String())

	log.Println(in.Name)
	log.Println(string(in.Compose))

	log.Println("[cluster] Success: created cluster")
	return &CreateReply{}, nil
}

// List implements cluster.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Println("[cluster] List", in.String())

	log.Println("[cluster] Success: list")
	return &ListReply{}, nil
}

// Status implements cluster.Server
func (s *Server) Status(ctx context.Context, in *StatusRequest) (*StatusReply, error) {
	log.Println("[cluster] Status", in.String())

	log.Println("[cluster] Success: list")
	return &StatusReply{}, nil
}

// Update implements cluster.Server
func (s *Server) Update(ctx context.Context, in *UpdateRequest) (*UpdateReply, error) {
	log.Println("[cluster] Update", in.String())

	log.Println("[cluster] Success: list")
	return &UpdateReply{}, nil
}

// Remove implements cluster.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	log.Println("[cluster] Remove", in.String())

	log.Println("[cluster] Success: removed", in.Id)
	return &RemoveReply{}, nil
}
