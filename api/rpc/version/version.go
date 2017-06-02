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

// ValidateGtwURL validation of the gtw url
func (s *Server) ValidateGtwURL(ctx context.Context, in *ValidateGtwURLRequest) (*ValidateGtwURLReply, error) {
	return &ValidateGtwURLReply{
		Reply: "Success, click back on our browser to go back to the login page",
	}, nil
}
