package topic

import (
	"fmt"
	"github.com/appcelerator/amp-kafka-pilot/pilot/api/admin"
	"github.com/appcelerator/amp/data/kafka"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"path"
)

const (
	topicsRootKey = "topics"
)

// Server is used to implement topic.TopicServer
type Server struct {
	Store storage.Interface
	Kafka kafka.Kafka
}

// Create implements topic.TopicServer
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	request := &admin.CreateTopicRequest{
		Topic: &admin.TopicEntry{
			Name:              in.Topic.Name,
			Partitions:        in.Topic.Partitions,
			ReplicationFactor: in.Topic.ReplicationFactor,
		},
	}
	created, err := s.Kafka.Admin().CreateTopic(ctx, request)
	if err != nil {
		return nil, err
	}
	topic := &TopicEntry{
		Id:                stringid.GenerateNonCryptoID(),
		Name:              created.Topic.Name,
		Partitions:        created.Topic.Partitions,
		ReplicationFactor: created.Topic.ReplicationFactor,
	}
	if err := s.Store.Create(ctx, path.Join(topicsRootKey, topic.Id), topic, nil, 0); err != nil {
		s.Kafka.Admin().DeleteTopic(ctx, &admin.DeleteTopicRequest{TopicName: topic.Name})
		return nil, err
	}
	return &CreateReply{Topic: topic}, nil
}

// List implements topic.TopicServer
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	var topics []proto.Message
	if err := s.Store.List(ctx, topicsRootKey, storage.Everything, &TopicEntry{}, &topics); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("Topic id not found: %s\n", in.Id)
	}

	request := &admin.DeleteTopicRequest{TopicName: topic.Name}
	if _, err := s.Kafka.Admin().DeleteTopic(ctx, request); err != nil {
		return nil, err
	}

	if err := s.Store.Delete(ctx, path.Join(topicsRootKey, in.Id), false, nil); err != nil {
		return nil, err
	}

	return &DeleteReply{Id: in.Id}, nil
}
