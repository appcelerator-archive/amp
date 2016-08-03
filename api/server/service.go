package server

import (
	"github.com/appcelerator/amp/api/rpc/service"
	"golang.org/x/net/context"
)

// serviceService is used to implement service.ServiceServer
type serviceService struct {
}

// Create implements service.ServiceServer
func (s *serviceService) Create(ctx context.Context, in *service.CreateRequest) (*service.CreateReply, error) {
	return nil, nil
}
