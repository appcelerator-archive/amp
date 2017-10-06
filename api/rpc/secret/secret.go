package secret

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
	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{Name: in.Name},
		Data:        in.Data,
	}
	reply, err := s.Docker.GetClient().SecretCreate(ctx, spec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Infoln("Successfully created secret:", in.Name)
	return &CreateReply{Id: reply.ID}, nil
}

func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	options := types.SecretListOptions{}
	secrets, err := s.Docker.GetClient().SecretList(ctx, options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Infoln("Successfully listed secret")
	reply := &ListReply{}
	for _, secret := range secrets {
		reply.Entries = append(reply.Entries, &SecretEntry{Id: secret.ID, Name: secret.Spec.Name})
	}
	return reply, nil
}

func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	if err := s.Docker.GetClient().SecretRemove(ctx, in.Id); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Infoln("Successfully removed secret:", in.Id)
	return &RemoveReply{Id: in.Id}, nil
}
