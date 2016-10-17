package logs_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/server"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"math/rand"
	"strconv"
)

const (
	natsClusterID          = "test-cluster"
	natsClientID           = "amplifier-log-test"
	defaultNumberOfEntries = 50

	testMessage     = "test message "
	testServiceName = "testServiceName"
	testStackName   = "testStackName"
)

var (
	testContainerId = stringid.GenerateNonCryptoID()
	testNodeId      = stringid.GenerateNonCryptoID()
	testServiceId   = stringid.GenerateNonCryptoID()
	testStackId     = stringid.GenerateNonCryptoID()
	ctx             context.Context
	client          LogsClient
	sc              stan.Conn
)

func TestMain(m *testing.M) {
	config, conn := server.StartTestServer()
	ctx = context.Background()
	client = NewLogsClient(conn)

	var err error
	sc, err = stan.Connect(natsClusterID, natsClientID+strconv.Itoa(rand.Int()), stan.NatsURL(config.NatsURL), stan.ConnectWait(60*time.Second))
	if err != nil {
		fmt.Println("failed to connect to nats")
		os.Exit(1)
	}
	defer func() {
		sc.Close()
	}()

	produceLogEntries(100)

	// Wait for entries to be indexed
	for {
		time.Sleep(1 * time.Second)
		r, err := client.Get(ctx, &GetRequest{Service: testServiceId})
		if err != nil {
			continue
		}
		if len(r.Entries) == 100 {
			break
		}
	}

	os.Exit(m.Run())
}

func TestShouldGetAHundredLogEntriesByDefault(t *testing.T) {
	expected := 100
	actual := -1
	for i := 0; i < 60; i++ {
		r, err := client.Get(ctx, &GetRequest{})
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
	r, err := client.Get(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomContainerId := r.Entries[0].ContainerId

	// Then filter by this container id
	r, err = client.Get(ctx, &GetRequest{Container: randomContainerId})
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
	r, err := client.Get(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomNodeId := r.Entries[0].NodeId

	// Then filter by this node id
	r, err = client.Get(ctx, &GetRequest{Node: randomNodeId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomNodeId, entry.NodeId)
	}
}

func TestShouldFilterByService(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Service: testServiceId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceId) || strings.HasPrefix(entry.ServiceId, testServiceId))
	}

	r, err = client.Get(ctx, &GetRequest{Service: testServiceName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceName))
	}
}

func TestShouldFilterByMessage(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Message: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Contains(t, strings.ToLower(entry.Message), "test")
	}
}

func TestShouldFilterByStack(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Stack: testStackId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackId) || strings.HasPrefix(entry.StackId, testStackId))
	}

	r, err = client.Get(ctx, &GetRequest{Stack: testStackName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackName) || strings.HasPrefix(entry.StackId, testStackName))
	}
}

func TestShouldFetchGivenNumberOfEntries(t *testing.T) {
	for i := int64(1); i < 100; i += 10 {
		r, err := client.Get(ctx, &GetRequest{Size: i})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, i, int64(len(r.Entries)))
	}
}

func produceLogEntries(howMany int) error {
	for i := 0; i < howMany; i++ {
		message, err := proto.Marshal(&LogEntry{
			ContainerId: testContainerId,
			Message:     testMessage + strconv.Itoa(rand.Int()),
			NodeId:      testNodeId,
			ServiceId:   testServiceId,
			ServiceName: testServiceName,
			StackId:     testStackId,
			StackName:   testStackName,
			Timestamp:   time.Now().Format(time.RFC3339Nano),
			TimeId:      time.Now().Format(time.RFC3339Nano),
		})
		err = sc.Publish(NatsLogTopic, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func listenToLogEntries(stream Logs_GetStreamClient, howMany int) (chan *LogEntry, error) {
	entries := make(chan *LogEntry, howMany)
	entryCount := 0
	timeout := time.After(5 * time.Second)

	defer func() {
		close(entries)
	}()

	for {
		entry, err := stream.Recv()
		select {
		case entries <- entry:
			if err != nil {
				return nil, err
			}
			entryCount++
			if entryCount == howMany {
				return entries, nil
			}
		case <-timeout:
			return entries, nil
		}
	}
}

func TestShouldStreamLogs(t *testing.T) {
	stream, err := client.GetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
}

func TestShouldStreamAndFilterByContainer(t *testing.T) {
	stream, err := client.GetStream(ctx, &GetRequest{Container: testContainerId})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testContainerId, entry.ContainerId)
	}
}

func TestShouldStreamAndFilterByNode(t *testing.T) {
	stream, err := client.GetStream(ctx, &GetRequest{Node: testNodeId})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testNodeId, entry.NodeId)
	}
}

func TestShouldStreamAndFilterByService(t *testing.T) {
	stream, err := client.GetStream(ctx, &GetRequest{Service: "testService"})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, "testService") || strings.HasPrefix(entry.ServiceId, "testService"))
	}
}

func TestShouldStreamAndFilterByMessage(t *testing.T) {
	stream, err := client.GetStream(ctx, &GetRequest{Message: testMessage})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Message), testMessage)
	}
}

func TestShouldStreamAndFilterCaseInsensitivelyByMessage(t *testing.T) {
	stream, err := client.GetStream(ctx, &GetRequest{Message: strings.ToUpper(testMessage)})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Message), testMessage)
	}
}

func TestShouldStreamAndFilterByStack(t *testing.T) {
	stream, err := client.GetStream(ctx, &GetRequest{Stack: "testStack"})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.StackName, "testStack") || strings.HasPrefix(entry.StackId, "testStack"))
	}
}
