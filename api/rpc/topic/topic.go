package topic

import (
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"path"
	"strings"
)

const (
	topicsRootKey = "topics"
)

// Server is used to implement topic.TopicServer
type Server struct {
	Store         storage.Interface
	NatsStreaming ns.NatsStreaming
}

// Create implements topic.TopicServer
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	reply, err := s.List(ctx, &ListRequest{})
	if err != nil {
		return nil, err
	}
	for _, topic := range reply.Topics {
		if strings.EqualFold(topic.Name, in.Topic.Name) {
			return nil, grpc.Errorf(codes.AlreadyExists, "Topic already exists: %s", in.Topic.Name)
		}
	}
	topic := &TopicEntry{
		Id:   stringid.GenerateNonCryptoID(),
		Name: in.Topic.Name,
	}
	if err := s.Store.Create(ctx, path.Join(topicsRootKey, topic.Id), topic, nil, 0); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return &CreateReply{Topic: topic}, nil
}

// List implements topic.TopicServer
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	var topics []proto.Message
	if err := s.Store.List(ctx, topicsRootKey, storage.Everything, &TopicEntry{}, &topics); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	reply := &ListReply{}
	for _, topic := range topics {
		reply.Topics = append(reply.Topics, topic.(*TopicEntry))
	}
	return reply, nil
}

// Delete implements topic.TopicServer
func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*DeleteReply, error) {
	topic := &TopicEntry{}
	if err := s.Store.Get(ctx, path.Join(topicsRootKey, in.Id), topic, false); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Topic not found: %s", in.Id)
	}

	if err := s.Store.Delete(ctx, path.Join(topicsRootKey, in.Id), false, nil); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}

	return &DeleteReply{Topic: topic}, nil
}
