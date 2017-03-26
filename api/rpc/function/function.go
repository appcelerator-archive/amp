package function

import (
	"log"

	"github.com/appcelerator/amp/data/functions"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server is used to implement function.FunctionServer
type Server struct {
	Functions     functions.Interface
	NatsStreaming ns.NatsStreaming
}

func convertError(err error) error {
	switch err {
	case functions.InvalidName:
	case functions.InvalidImage:
		return grpc.Errorf(codes.InvalidArgument, err.Error())
	case functions.FunctionAlreadyExists:
		return grpc.Errorf(codes.AlreadyExists, err.Error())
	}
	return grpc.Errorf(codes.Internal, err.Error())
}

// Create implements function.Server
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	log.Println("rpc-function: Create", in.String())
	function, err := s.Functions.CreateFunction(ctx, in.Name, in.Image)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully created function", function.String())
	return &CreateReply{Function: function}, nil
}

// List implements function.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	functions, err := s.Functions.ListFunctions(ctx)
	if err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully list functions")
	return &ListReply{Functions: functions}, nil
}

// Delete implements function.Server
func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*empty.Empty, error) {
	log.Println("rpc-function: Delete", in.String())

	// Check if function exists
	function, err := s.Functions.GetFunction(ctx, in.Id)
	if err != nil {
		return nil, convertError(err)
	}
	if function == nil {
		return nil, grpc.Errorf(codes.NotFound, "function not found: %s", in.Id)
	}

	if err := s.Functions.DeleteFunction(ctx, in.Id); err != nil {
		return nil, convertError(err)
	}
	log.Println("Successfully deleted function", in.Id)
	return &empty.Empty{}, nil
}
