package stack

import (
	"github.com/appcelerator/amp/data/storage"
	"golang.org/x/net/context"
)

// Server is used to implement stack.StackService
type Server struct {
	Store storage.Interface
}

// Create implements stack.ServerService Create
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	stack, err := NewStackfromYaml(ctx, in.StackDefinition)
	if err != nil {
		return nil, err
	}
	stack.Name = in.StackName
	reply := CreateReply{
		StackId: stack.Id,
	}
	return &reply, nil
}

// Up implements stack.ServerService Up
func (s *Server) Up(ctx context.Context, in *UpRequest) (*UpReply, error) {
	stack, err := NewStackfromYaml(ctx, in.Stackfile)
	if err != nil {
		return nil, err
	}
	stack.Name = in.StackName
	reply := UpReply{
		StackId: stack.Id,
	}
	return &reply, nil
}
