package api

import (
	"github.com/appcelerator/amp-docker-kafka/pilot/api/admin"
	"google.golang.org/grpc"
)

// NewAdminClient create a new administration client
func NewAdminClient(host string) (admin.AdminClient, error) {
	clientConn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return admin.NewAdminClient(clientConn), nil
}
