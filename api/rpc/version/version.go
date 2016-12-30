package version

import (
	"golang.org/x/net/context"
)

type Details struct {
	Version    string
	Build      string
	ConfigAddr string
	Port       string
	GoVersion  string
	Os         string
	Arch       string
}

type VersionConfig struct {
	AMP       *Details
	Amplifier *Details
}

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

// Returns Amplifier info from Server.go
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

// Checks if AMP is connected to Amplifier
func (v VersionConfig) AmplifierOK() bool {
	return v.Amplifier != nil
}
