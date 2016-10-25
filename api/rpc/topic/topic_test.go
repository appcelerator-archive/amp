package topic_test

import (
	. "github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/api/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"os"
	"testing"
)

var (
	ctx    context.Context
	client TopicClient
)

func TestMain(m *testing.M) {
	_, conn := server.StartTestServer()
	client = NewTopicClient(conn)
	ctx = context.Background()
	os.Exit(m.Run())
}

func TestShouldCreateAndDeleteATopic(t *testing.T) {
	created, err := client.Create(ctx, &CreateRequest{
		Topic: &TopicEntry{
			Name: "test-topic",
		},
	})
	assert.NoError(t, err)

	_, err = client.Delete(ctx, &DeleteRequest{
		Id: created.Topic.Id,
	})
	assert.NoError(t, err)
}

func TestShouldFailWhenCreatingAnAlreadyExistingTopic(t *testing.T) {
	created, err := client.Create(ctx, &CreateRequest{
		Topic: &TopicEntry{
			Name: "test-topic",
		},
	})
	assert.NoError(t, err)

	_, err = client.Create(ctx, &CreateRequest{
		Topic: &TopicEntry{
			Name: "test-topic",
		},
	})
	assert.Error(t, err)

	_, err = client.Delete(ctx, &DeleteRequest{
		Id: created.Topic.Id,
	})
	assert.NoError(t, err)
}

func TestShouldListCreatedTopics(t *testing.T) {
	r1, err := client.Create(ctx, &CreateRequest{Topic: &TopicEntry{Name: "test-topic-1"}})
	assert.NoError(t, err)
	r2, err := client.Create(ctx, &CreateRequest{Topic: &TopicEntry{Name: "test-topic-2"}})
	assert.NoError(t, err)
	r3, err := client.Create(ctx, &CreateRequest{Topic: &TopicEntry{Name: "test-topic-3"}})
	assert.NoError(t, err)

	reply, err := client.List(ctx, &ListRequest{})
	assert.NoError(t, err)
	assert.Contains(t, reply.Topics, r1.Topic)
	assert.Contains(t, reply.Topics, r2.Topic)
	assert.Contains(t, reply.Topics, r3.Topic)

	_, err = client.Delete(ctx, &DeleteRequest{Id: r1.Topic.Id})
	assert.NoError(t, err)
	_, err = client.Delete(ctx, &DeleteRequest{Id: r2.Topic.Id})
	assert.NoError(t, err)
	_, err = client.Delete(ctx, &DeleteRequest{Id: r3.Topic.Id})
	assert.NoError(t, err)

	reply, err = client.List(ctx, &ListRequest{})
	assert.NoError(t, err)
	assert.NotContains(t, reply.Topics, r1.Topic)
	assert.NotContains(t, reply.Topics, r2.Topic)
	assert.NotContains(t, reply.Topics, r3.Topic)
}
