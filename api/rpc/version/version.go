package version

import (
	"golang.org/x/net/context"
)

//Details detailled information
type Details struct {
	Version    string
	Build      string
	ConfigAddr string
	Port       string
	GoVersion  string
	Os         string
	Arch       string
}

// Config version information
type Config struct {
	AMP       *Details
	Amplifier *Details
}

// Server server information
type Server struct {
	Version   string
	Port      string
	GoVersion string
	Os        string
	Arch      string
}

func init() {
	// no setup
}

// List Returns Amplifier info from Server.go
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {

	response := &ListReply{
		Reply: &VersionInfo{
			Version:   s.Version,
			Port:      s.Port,
			Goversion: s.GoVersion,
			Os:        s.Os,
			Arch:      s.Arch,
		},
	}
	return response, nil
}

// AmplifierOK Checks if AMP is connected to Amplifier
func (v Config) AmplifierOK() bool {
	return v.Amplifier != nil
}
