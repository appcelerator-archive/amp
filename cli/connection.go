package cli

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewClientConn is a helper function that wraps the steps involved in setting up a grpc client connection to the API.
func NewClientConn(addr string, token string, skipVerify bool) (*grpc.ClientConn, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: skipVerify}
	creds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
		grpc.WithPerRPCCredentials(&LoginCredentials{Token: token}),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
	}
	return grpc.Dial(addr, opts...)
}
