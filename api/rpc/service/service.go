package service

import (
	"log"

	"github.com/docker/docker/client"
	//	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

var (
	// https://docs.docker.com/engine/reference/api/docker_remote_api/
	// `docker version` -> Server API version  => Docker 1.12x
	defaultVersion = "v1.25"
	defaultHeaders = map[string]string{"User-Agent": "engine-api-cli-1.0"}
	dockerSock     = "unix:///var/run/docker.sock"
	docker         *client.Client
	err            error
)

// Service is used to implement ServiceServer
type Service struct{}

func init() {
	docker, err = client.NewClient(dockerSock, defaultVersion, nil, defaultHeaders)
	if err != nil {
		// fail fast
		panic(err)
	}
}

// Create implements ServiceServer
func (s *Service) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	log.Println(in)
	reply := &CreateReply{}
	reply.Message = in.String()
	return reply, nil
}
