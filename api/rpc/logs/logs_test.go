package logs_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/server"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"math/rand"
	"strconv"
)

const (
	defaultNumberOfEntries = 50
	testContainerId        = "testContainerId"
	testMessage            = "test message "
	testNodeId             = "testNodeId"
	testServiceId          = "testServiceId"
	testServiceName        = "testServiceName"
	testStackId            = "testStackId"
	testStackName          = "testStackName"
)

var (
	ctx      context.Context
	client   logs.LogsClient
	producer sarama.SyncProducer
)

func TestMain(m *testing.M) {
	config, conn := server.StartTestServer()
	client = logs.NewLogsClient(conn)
	ctx = context.Background()

	var err error
	producer, err = sarama.NewSyncProducer([]string{config.KafkaURL}, nil)
	if err != nil {
		fmt.Println("Cannot create kafka producer")
		os.Exit(1)
	}
	defer func() {
		producer.Close()
	}()

	produceLogEntries(1000)

	os.Exit(m.Run())
}

func TestShouldGetAHundredLogEntriesByDefault(t *testing.T) {
	expected := 100
	actual := -1
	for i := 0; i < 60; i++ {
		r, err := client.Get(ctx, &logs.GetRequest{})
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		actual = len(r.Entries)
		if actual == expected {
			break
		}
		time.Sleep(1 * time.Second)
	}
	assert.Equal(t, expected, actual)
}

func TestShouldFilterByContainer(t *testing.T) {
	// First, get a random container id
	r, err := client.Get(ctx, &logs.GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomContainerId := r.Entries[0].ContainerId

	// Then filter by this container id
	r, err = client.Get(ctx, &logs.GetRequest{Container: randomContainerId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomContainerId, entry.ContainerId)
	}
}

func TestShouldFilterByNode(t *testing.T) {
	// First, get a random node id
	r, err := client.Get(ctx, &logs.GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomNodeId := r.Entries[0].NodeId

	// Then filter by this node id
	r, err = client.Get(ctx, &logs.GetRequest{Node: randomNodeId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomNodeId, entry.NodeId)
	}
}

func TestShouldFilterByService(t *testing.T) {
	r, err := client.Get(ctx, &logs.GetRequest{Service: "et"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, "et") || strings.HasPrefix(entry.ServiceId, "et"))
	}
}

func TestShouldFilterByMessage(t *testing.T) {
	r, err := client.Get(ctx, &logs.GetRequest{Message: "info"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Contains(t, strings.ToLower(entry.Message), "info")
	}
}

/*
func TestShouldFilterByStack(t *testing.T) {
	r, err := client.Get(ctx, &logs.GetRequest{Stack: "testStack"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, "testStack") || strings.HasPrefix(entry.StackId, "testStack"))
	}
}
*/
func TestShouldFetchGivenNumberOfEntries(t *testing.T) {
	for i := int64(1); i < 200; i += 10 {
		r, err := client.Get(ctx, &logs.GetRequest{Size: i})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, i, int64(len(r.Entries)))
	}
}

func produceLogEntries(howMany int) error {
	for i := 0; i < howMany; i++ {
		message, err := proto.Marshal(&logs.LogEntry{
			ContainerId: testContainerId,
			Message:     testMessage + strconv.Itoa(rand.Int()),
			NodeId:      testNodeId,
			ServiceId:   testServiceId,
			ServiceName: testServiceName,
			StackId:     testStackId,
			StackName:   testStackName,
			Timestamp:   time.Now().Format("2006-01-02T15:04:05.999"),
			TimeId:      time.Now().Format(time.RFC3339Nano),
		})
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "amp-logs",
			Value: sarama.ByteEncoder(message),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func listenToLogEntries(t *testing.T, stream logs.Logs_GetStreamClient, howMany int) chan *logs.LogEntry {
	entries := make(chan *logs.LogEntry, howMany)
	entryCount := 0
	timeout := time.After(60 * time.Second)

	defer func() {
		close(entries)
	}()

	for {
		entry, err := stream.Recv()
		select {
		case entries <- entry:
			if err != nil {
				t.Error(err)
			}
			entryCount++
			if entryCount == howMany {
				return entries
			}
		case <-timeout:
			return entries
		}
	}
}

func TestShouldStreamLogs(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
}

func TestShouldStreamAndFilterByContainer(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{Container: testContainerId})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testContainerId, entry.ContainerId)
	}
}

func TestShouldStreamAndFilterByNode(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{Node: testNodeId})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testNodeId, entry.NodeId)
	}
}

func TestShouldStreamAndFilterByService(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{Service: "testService"})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, "testService") || strings.HasPrefix(entry.ServiceId, "testService"))
	}
}

func TestShouldStreamAndFilterByMessage(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{Message: testMessage})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Message), testMessage)
	}
}

func TestShouldStreamAndFilterCaseInsensitivelyByMessage(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{Message: strings.ToUpper(testMessage)})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Message), testMessage)
	}
}

func TestShouldStreamAndFilterByStack(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{Stack: "testStack"})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.StackName, "testStack") || strings.HasPrefix(entry.StackId, "testStack"))
	}
}
