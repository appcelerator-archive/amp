package config

import (
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/pkg/docker"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server server information
type Server struct {
	Docker *docker.Docker
}

func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	spec := swarm.ConfigSpec{
		Annotations: swarm.Annotations{Name: in.Name},
		Data:        in.Data,
	}
	reply, err := s.Docker.GetClient().ConfigCreate(ctx, spec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Infoln("Successfully created config:", in.Name)
	return &CreateReply{Id: reply.ID}, nil
}

func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	options := types.ConfigListOptions{}
	configs, err := s.Docker.GetClient().ConfigList(ctx, options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Infoln("Successfully listed config")
	reply := &ListReply{}
	for _, config := range configs {
		reply.Entries = append(reply.Entries, &ConfigEntry{Id: config.ID, Name: config.Spec.Name})
	}
	return reply, nil
}

func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	if err := s.Docker.GetClient().ConfigRemove(ctx, in.Id); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Infoln("Successfully removed config:", in.Id)
	return &RemoveReply{Id: in.Id}, nil
}
