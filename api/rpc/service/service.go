package service

import (
	"golang.org/x/net/context"
)

// Service is used to implement ServiceServer
type Service struct {
}

// Create implements ServiceServer
func (s *Service) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	return nil, nil
}
