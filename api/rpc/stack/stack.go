package stack

import (
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"golang.org/x/net/context"
)

// Server is used to implement stack.StackService
type Server struct {
	Store storage.Interface
}

// Create implements stack.ServerService Create
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	stack, err := parseStackYaml(in.StackDefinition)
	if err != nil {
		return nil, err
	}
	stackID := stringid.GenerateNonCryptoID()
	stack.Id = stackID
	s.Store.Create(ctx, "stacks/"+stackID, stack, nil, 0)
	reply := CreateReply{
		StackId: stringid.GenerateNonCryptoID(),
	}
	return &reply, nil
}

// Up implements stack.ServerService Up
func (s *Server) Up(ctx context.Context, in *UpRequest) (*UpReply, error) {
	stack, err := parseStackYaml(in.Stackfile)
	if err != nil {
		return nil, err
	}
	stackID := stringid.GenerateNonCryptoID()
	stack.Id = stackID
	s.Store.Create(ctx, "stacks/"+stackID, stack, nil, 0)
	reply := UpReply{
		StackId: stackID,
	}
	return &reply, nil
}
