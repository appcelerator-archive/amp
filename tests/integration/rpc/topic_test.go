package tests

import (
	. "github.com/appcelerator/amp/api/rpc/topic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopicShouldCreateAndDeleteATopic(t *testing.T) {
	created, err := topicClient.Create(ctx, &CreateRequest{
		Topic: &TopicEntry{
			Name: "test-topic",
		},
	})
	assert.NoError(t, err)

	_, err = topicClient.Delete(ctx, &DeleteRequest{
		Id: created.Topic.Id,
	})
	assert.NoError(t, err)
}

func TestTopicShouldFailWhenCreatingAnAlreadyExistingTopic(t *testing.T) {
	created, err := topicClient.Create(ctx, &CreateRequest{
		Topic: &TopicEntry{
			Name: "test-topic",
		},
	})
	assert.NoError(t, err)

	_, err = topicClient.Create(ctx, &CreateRequest{
		Topic: &TopicEntry{
			Name: "test-topic",
		},
	})
	assert.Error(t, err)

	_, err = topicClient.Delete(ctx, &DeleteRequest{
		Id: created.Topic.Id,
	})
	assert.NoError(t, err)
}

func TestTopicShouldListCreatedTopics(t *testing.T) {
	r1, err := topicClient.Create(ctx, &CreateRequest{Topic: &TopicEntry{Name: "test-topic-1"}})
	assert.NoError(t, err)
	r2, err := topicClient.Create(ctx, &CreateRequest{Topic: &TopicEntry{Name: "test-topic-2"}})
	assert.NoError(t, err)
	r3, err := topicClient.Create(ctx, &CreateRequest{Topic: &TopicEntry{Name: "test-topic-3"}})
	assert.NoError(t, err)

	reply, err := topicClient.List(ctx, &ListRequest{})
	assert.NoError(t, err)
	assert.Contains(t, reply.Topics, r1.Topic)
	assert.Contains(t, reply.Topics, r2.Topic)
	assert.Contains(t, reply.Topics, r3.Topic)

	_, err = topicClient.Delete(ctx, &DeleteRequest{Id: r1.Topic.Id})
	assert.NoError(t, err)
	_, err = topicClient.Delete(ctx, &DeleteRequest{Id: r2.Topic.Id})
	assert.NoError(t, err)
	_, err = topicClient.Delete(ctx, &DeleteRequest{Id: r3.Topic.Id})
	assert.NoError(t, err)

	reply, err = topicClient.List(ctx, &ListRequest{})
	assert.NoError(t, err)
	assert.NotContains(t, reply.Topics, r1.Topic)
	assert.NotContains(t, reply.Topics, r2.Topic)
	assert.NotContains(t, reply.Topics, r3.Topic)
}
