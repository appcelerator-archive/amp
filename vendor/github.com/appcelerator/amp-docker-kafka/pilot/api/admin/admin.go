package admin

import (
	"github.com/appcelerator/amp-docker-kafka/pilot"
	"golang.org/x/net/context"
)

// Server implement topic.AdminServer
type Server struct {
	Pilot *pilot.KafkaPilot
}

// CreateTopic create a topic
func (s *Server) CreateTopic(ctx context.Context, in *CreateTopicRequest) (*CreateTopicReply, error) {
	if err := s.Pilot.CreateTopic(in.Topic.Name, in.Topic.Partitions, in.Topic.ReplicationFactor); err != nil {
		return nil, err
	}
	return &CreateTopicReply{
		Topic: &TopicEntry{
			Name:              in.Topic.Name,
			Partitions:        in.Topic.Partitions,
			ReplicationFactor: in.Topic.ReplicationFactor,
		},
	}, nil
}

// DeleteTopic delete a topic
func (s *Server) DeleteTopic(ctx context.Context, in *DeleteTopicRequest) (*DeleteTopicReply, error) {
	if err := s.Pilot.DeleteTopic(in.TopicName); err != nil {
		return nil, err
	}
	return &DeleteTopicReply{}, nil
}
