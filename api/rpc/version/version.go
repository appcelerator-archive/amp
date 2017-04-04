package version

import (
	"golang.org/x/net/context"
)

// Server server information
type Server struct {
	Info *Info
}

// Get Returns Amplifier info from Server.go
func (s *Server) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	return &GetReply{Info: s.Info}, nil
}
