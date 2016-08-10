package server

import (
	"github.com/appcelerator/amp/api/rpc/logs"
	"golang.org/x/net/context"
)

// logService is used to implement log.LogServer
type logsService struct {
}

// Get implements log.LogServer
func (s *logsService) Get(ctx context.Context, in *logs.GetRequest) (*logs.GetReply, error) {
	return &logs.GetReply{}, nil
}
