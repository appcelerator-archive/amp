package tests

import (
	"log"
	"strings"
	"testing"
	"time"

	. "github.com/appcelerator/amp/api/rpc/logs"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats"
	"github.com/stretchr/testify/assert"
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
	sc              stan.Conn
	err             error
)

func TestLogsInit(t *testing.T) {
	var err error
	log.Printf("Connecting to nats: %s\n", config.NatsURL)
	nc, err := nats.Connect(config.NatsURL, nats.Timeout(60*time.Second))
	if err != nil {
		t.Errorf("Unable to connect to NATS on: %s\n%v", config.NatsURL, err)
		return
	}
	//runtime.Nats, err = stan.Connect(natsClusterID, natsClientID+strconv.Itoa(rand.Int()), stan.NatsConn(nc), stan.ConnectWait(defaultTimeOut))
	sc, err = stan.Connect(natsClusterID, natsClientID+strconv.Itoa(rand.Int()), stan.NatsConn(nc), stan.ConnectWait(60*time.Second))
	if err != nil {
		t.Errorf("failed to connect to nats: %v\n", err)
		return
	}
	log.Println("Connected to nats")
	err = produceLogEntries(110)
	if err != nil {
		t.Errorf("failed to produce log entries: %v\n", err)
		return
	}
	// Wait for entries to be indexed
	for {
		time.Sleep(1 * time.Second)
		r, err := logsClient.Get(ctx, &GetRequest{Service: testServiceId})
		if err != nil {
			continue
		}
		if len(r.Entries) == 100 {
			break
		}
	}
}

func TestLogsShouldGetAHundredLogEntriesByDefault(t *testing.T) {
	expected := 100
	actual := -1
	for i := 0; i < 60; i++ {
		r, err := logsClient.Get(ctx, &GetRequest{})
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

func TestLogsShouldFilterByContainer(t *testing.T) {
	// First, get a random container id
	r, err := logsClient.Get(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomContainerId := r.Entries[0].ContainerId

	// Then filter by this container id
	r, err = logsClient.Get(ctx, &GetRequest{Container: randomContainerId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomContainerId, entry.ContainerId)
	}
}

func TestLogsShouldFilterByNode(t *testing.T) {
	// First, get a random node id
	r, err := logsClient.Get(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomNodeId := r.Entries[0].NodeId

	// Then filter by this node id
	r, err = logsClient.Get(ctx, &GetRequest{Node: randomNodeId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomNodeId, entry.NodeId)
	}
}

func TestLogsShouldFilterByService(t *testing.T) {
	r, err := logsClient.Get(ctx, &GetRequest{Service: testServiceId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceId) || strings.HasPrefix(entry.ServiceId, testServiceId))
	}

	r, err = logsClient.Get(ctx, &GetRequest{Service: testServiceName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceName))
	}
}

func TestLogsShouldFilterByMessage(t *testing.T) {
	r, err := logsClient.Get(ctx, &GetRequest{Message: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Contains(t, strings.ToLower(entry.Message), "test")
	}
}

func TestLogsShouldFilterByStack(t *testing.T) {
	r, err := logsClient.Get(ctx, &GetRequest{Stack: testStackId})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackId) || strings.HasPrefix(entry.StackId, testStackId))
	}

	r, err = logsClient.Get(ctx, &GetRequest{Stack: testStackName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackName) || strings.HasPrefix(entry.StackId, testStackName))
	}
}

func TestLogsShouldFetchGivenNumberOfEntries(t *testing.T) {
	for i := int64(1); i < 100; i += 10 {
		r, err := logsClient.Get(ctx, &GetRequest{Size: i})
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

func TestLogsShouldStreamLogs(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
}

func TestLogsShouldStreamAndFilterByContainer(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Container: testContainerId})
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

func TestLogsShouldStreamAndFilterByNode(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Node: testNodeId})
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

func TestLogsShouldStreamAndFilterByService(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Service: "testService"})
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

func TestogsShouldStreamAndFilterByMessage(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Message: testMessage})
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

func TestLogsShouldStreamAndFilterCaseInsensitivelyByMessage(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Message: strings.ToUpper(testMessage)})
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

func TestLogsShouldStreamAndFilterByStack(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Stack: "testStack"})
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

func TestLogsEnd(t *testing.T) {
	sc.Close()
}
