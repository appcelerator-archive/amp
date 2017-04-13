package stack

import (
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// Server is used to implement stack.StackServer
type Server struct {
}

// Deploy implements stack.Server
func (s *Server) Deploy(ctx context.Context, in *DeployRequest) (*DeployReply, error) {
	log.Println("[stack] Deploy", in.String())

	log.Println(in.Name)
	log.Println(string(in.Compose))

	log.Println("[stack] Success: created stack")
	return &DeployReply{}, nil
}

// List implements stack.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Println("[stack] List", in.String())

	log.Println("[stack] Success: list")
	return &ListReply{}, nil
}

// Remove implements stack.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*empty.Empty, error) {
	log.Println("[stack] Remove", in.String())

	log.Println("[stack] Success: removed", in.Id)
	return &empty.Empty{}, nil
}
