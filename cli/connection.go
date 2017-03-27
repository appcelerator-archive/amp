package cli

import (
	"time"

	"google.golang.org/grpc"
)

// NewClientConn is a helper function that wraps the steps involved in setting up a grpc client connection to the API.
func NewClientConn(addr string, token string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
		//grpc.WithPerRPCCredentials(token),
	)
}
