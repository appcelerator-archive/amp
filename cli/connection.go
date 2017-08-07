package cli

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	shortConnectionTimeout = time.Second
	longConnectionTimeout  = 10 * time.Second
)

// if the address is a local one, return a short timeout, else return a longer timeout
func getTimeout(addr string) (time.Duration, error) {
	// a parsable URL needs a scheme, just add it if it's not there
	if !strings.Contains(addr, "://") {
		addr = fmt.Sprintf("scheme://%s", addr)
	}
	u, err := url.Parse(addr)
	if err != nil {
		return 0, err
	}
	host := u.Hostname()
	switch host {
	case "localhost", "127.0.0.1", "::1":
		return shortConnectionTimeout, nil
	default:
		return longConnectionTimeout, nil
	}
}

// NewClientConn is a helper function that wraps the steps involved in setting up a grpc client connection to the API.
func NewClientConn(addr string, token string, skipVerify bool) (*grpc.ClientConn, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: skipVerify}
	creds := credentials.NewTLS(tlsConfig)
	connectionTimeout, err := getTimeout(addr)
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithTimeout(connectionTimeout),
		grpc.WithPerRPCCredentials(&LoginCredentials{Token: token}),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
	}
	return grpc.Dial(addr, opts...)
}
