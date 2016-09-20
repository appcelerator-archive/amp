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
	testServiceId          = "testServiceId"
	testServiceName        = "testServiceName"
	testNodeId             = "testNodeId"
	testContainerId        = "testContainerId"
	testMessage            = "test message "
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

func TestShouldFilterByContainerId(t *testing.T) {
	// First, get a random container id
	r, err := client.Get(ctx, &logs.GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomContainerId := r.Entries[0].ContainerId

	// Then filter by this container id
	r, err = client.Get(ctx, &logs.GetRequest{ContainerId: randomContainerId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomContainerId, entry.ContainerId)
	}
}

func TestShouldFilterByNodeId(t *testing.T) {
	// First, get a random node id
	r, err := client.Get(ctx, &logs.GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomNodeId := r.Entries[0].NodeId

	// Then filter by this node id
	r, err = client.Get(ctx, &logs.GetRequest{NodeId: randomNodeId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomNodeId, entry.NodeId)
	}
}

func TestShouldFilterByServiceId(t *testing.T) {
	// First, get a random service id
	r, err := client.Get(ctx, &logs.GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomServiceId := r.Entries[0].ServiceId

	// Then filter by this service id
	r, err = client.Get(ctx, &logs.GetRequest{ServiceId: randomServiceId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomServiceId, entry.ServiceId)
	}
}

func TestShouldFilterByServiceName(t *testing.T) {
	// First, get a random service name
	r, err := client.Get(ctx, &logs.GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomServiceName := r.Entries[0].ServiceName

	// Then filter by this service name
	r, err = client.Get(ctx, &logs.GetRequest{ServiceName: randomServiceName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomServiceName, entry.ServiceName)
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

func TestShouldFetchGivenNumberOfEntries(t *testing.T) {
	for i := int64(1); i < 200; i += 10 {
		r, err := client.Get(ctx, &logs.GetRequest{Size: i})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, i, int64(len(r.Entries)))
	}
}

func produceLogEntries(t *testing.T, howMany int) {
	for i := 0; i < howMany; i++ {
		message, err := proto.Marshal(&logs.LogEntry{
			Timestamp:   time.Now().Format(time.RFC3339Nano),
			TimeId:      time.Now().Format(time.RFC3339Nano),
			ServiceId:   testServiceId,
			ServiceName: testServiceName,
			NodeId:      testNodeId,
			ContainerId: testContainerId,
			Message:     testMessage + strconv.Itoa(rand.Int()),
		})
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "amp-logs",
			Value: sarama.ByteEncoder(message),
		})
		if err != nil {
			t.Error(err)
		}
	}
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
	go produceLogEntries(t, 100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
}

func TestShouldStreamAndFilterByContainerId(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{ContainerId: testContainerId})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(t, 100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testContainerId, entry.ContainerId)
	}
}

func TestShouldStreamAndFilterByNodeId(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{NodeId: testNodeId})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(t, 100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testNodeId, entry.NodeId)
	}
}

func TestShouldStreamAndFilterByServiceId(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{ServiceId: testServiceId})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(t, 100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testServiceId, entry.ServiceId)
	}
}

func TestShouldStreamAndFilterByServiceName(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{ServiceName: testServiceName})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(t, 100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testServiceName, entry.ServiceName)
	}
}

func TestShouldStreamAndFilterByMessage(t *testing.T) {
	stream, err := client.GetStream(ctx, &logs.GetRequest{Message: testMessage})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(t, 100)
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
	go produceLogEntries(t, 100)
	entries := listenToLogEntries(t, stream, defaultNumberOfEntries)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Message), testMessage)
	}
}
