package function

import (
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"path"
	"strings"
)

// Server is used to implement function.FunctionServer
type Server struct {
	Store         storage.Interface
	NatsStreaming ns.NatsStreaming
}

// GetOwners get organization owners
func (f *FunctionEntry) GetOwners() []string {
	return []string{f.Owner}
}

// Create implements function.Server
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	log.Println("rpc-function: Create", in.String())

	// Get requester
	requester, err := auth.GetRequester(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	// Validate the function
	fn := in.Function
	fn.Name = strings.TrimSpace(fn.Name)
	fn.Image = strings.TrimSpace(fn.Image)
	if fn.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "function name is mandatory")
	}
	if fn.Image == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "docker image is mandatory")
	}

	// Check if the function already exists
	reply, err := s.List(ctx, &ListRequest{})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "error listing functions: %v, err")
	}
	for _, fn := range reply.Functions {
		if strings.EqualFold(fn.Name, in.Function.Name) {
			return nil, grpc.Errorf(codes.AlreadyExists, "function already exists: %s", in.Function.Name)
		}
	}

	// Store the function
	fn.Id = stringid.GenerateNonCryptoID()
	fn.Owner = requester
	if err := s.Store.Create(ctx, path.Join(amp.EtcdFunctionRootKey, fn.Id), fn, nil, 0); err != nil {
		return nil, grpc.Errorf(codes.Internal, "error creating function: %v", err)
	}
	log.Println("Created function:", fn.String())
	return &CreateReply{Function: fn}, nil
}

// List implements function.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Println("rpc-function: List", in.String())
	var functions []proto.Message
	if err := s.Store.List(ctx, amp.EtcdFunctionRootKey, storage.Everything, &FunctionEntry{}, &functions); err != nil {
		return nil, grpc.Errorf(codes.Internal, "error listing functions: %v", err)
	}
	reply := &ListReply{}
	for _, function := range functions {
		reply.Functions = append(reply.Functions, function.(*FunctionEntry))
	}
	log.Println("Listed functions")
	return reply, nil
}

// Delete implements function.Server
func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*DeleteReply, error) {
	log.Println("rpc-function: Delete", in.String())

	function := &FunctionEntry{}
	if err := s.Store.Get(ctx, path.Join(amp.EtcdFunctionRootKey, in.Id), function, true); err != nil {
		return nil, grpc.Errorf(codes.Internal, "error deleting function: %s", err.Error())
	}
	// If there's no "id" in the answer, it means the function has not been found, so return nil
	if function.GetId() == "" {
		return nil, grpc.Errorf(codes.NotFound, "function not found: %s", in.Id)
	}

	// Check authorization
	requester, err := auth.GetRequester(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "user %s is not authorized to perform this action", requester)
	}
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  requester,
		Action:   auth.DeleteAction,
		Resource: auth.FunctionResource,
		Context: ladon.Context{
			"resource": function,
		},
	}); err != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "user %s is not authorized to perform this action", requester)
	}

	if err := s.Store.Delete(ctx, path.Join(amp.EtcdFunctionRootKey, in.Id), false, nil); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "function not found: %s", in.Id)
	}
	log.Println("Deleted function id: ", in.Id)
	return &DeleteReply{}, nil
}
