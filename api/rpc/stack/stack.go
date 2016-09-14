package stack

import (
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"golang.org/x/net/context"
)

// Stack is used to implement stack.StackServer
type Stack struct {
	Store storage.Interface
}

// Create implements stack.StackServer
func (stack *Stack) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	// Build reply
	reply := CreateReply{
		StackId: stringid.GenerateNonCryptoID(),
	}
	return &reply, nil
}

// Up implements stack.StackServer
func (stack *Stack) Up(ctx context.Context, in *UpRequest) (*UpReply, error) {
	reply := UpReply{
		StackId: stringid.GenerateNonCryptoID(),
	}
	return &reply, nil
}
