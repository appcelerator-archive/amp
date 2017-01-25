package tests

import (
	"log"
	"strings"
	"testing"
	"time"

	"math/rand"
	"strconv"

	. "github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/config"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/stretchr/testify/assert"
)

const (
	natsClientID           = "amplifier-log-test"
	defaultNumberOfEntries = 50

	testMessage     = "test message "
	testServiceName = "testservice"
	testStackName   = "teststack"
	testTaskName    = "test.0"
)

var (
	testContainerID = stringid.GenerateNonCryptoID()
	testNodeID      = stringid.GenerateNonCryptoID()
	testServiceID   = stringid.GenerateNonCryptoID()
	testStackID     = stringid.GenerateNonCryptoID()
	testTaskID      = stringid.GenerateNonCryptoID()
	sc              stan.Conn
)

func TestLogsInit(t *testing.T) {
	var err error
	log.Printf("Connecting to nats: %s\n", amp.NatsDefaultURL)
	nc, err := nats.Connect(amp.NatsDefaultURL, nats.Timeout(60*time.Second))
	if err != nil {
		t.Errorf("Unable to connect to NATS on: %s\n%v", amp.NatsDefaultURL, err)
		return
	}
	sc, err = stan.Connect(amp.NatsClusterID, natsClientID+strconv.Itoa(rand.Int()), stan.NatsConn(nc), stan.ConnectWait(60*time.Second))
	if err != nil {
		t.Errorf("failed to connect to nats: %v\n", err)
		return
	}
	log.Println("Connected to nats")
	err = produceLogEntries(100, false)
	if err != nil {
		t.Errorf("failed to produce log entries: %v\n", err)
		return
	}
	// Wait for entries to be indexed
	for {
		time.Sleep(1 * time.Second)
		r, err := logsClient.Get(ctx, &GetRequest{Service: testServiceID})
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
	randomContainerID := r.Entries[0].ContainerId

	// Then filter by this container id
	r, err = logsClient.Get(ctx, &GetRequest{Container: randomContainerID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomContainerID, entry.ContainerId)
	}
}

func TestLogsShouldFilterByNode(t *testing.T) {
	// First, get a random node id
	r, err := logsClient.Get(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	randomNodeID := r.Entries[0].NodeId

	// Then filter by this node id
	r, err = logsClient.Get(ctx, &GetRequest{Node: randomNodeID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, randomNodeID, entry.NodeId)
	}
}

func TestLogsShouldFilterByService(t *testing.T) {
	r, err := logsClient.Get(ctx, &GetRequest{Service: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceID))
	}

	r, err = logsClient.Get(ctx, &GetRequest{Service: testServiceName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceID))
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
	r, err := logsClient.Get(ctx, &GetRequest{Stack: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackName) || strings.HasPrefix(entry.StackId, testStackID))
	}

	r, err = logsClient.Get(ctx, &GetRequest{Stack: testStackName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackName) || strings.HasPrefix(entry.StackId, testStackID))
	}
}

func TestLogsShouldFilterByInfrastructureRole(t *testing.T) {
	r, err := logsClient.Get(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, entry.Role != amp.InfrastructureRole)
	}

	r, err = logsClient.Get(ctx, &GetRequest{Service: "elastic", Infra: true})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	gotInfraEntry := false
	for _, entry := range r.Entries {
		if entry.Role == amp.InfrastructureRole {
			gotInfraEntry = true
			break
		}
	}
	assert.True(t, gotInfraEntry)
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

func TestLogsShouldStreamLogs(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, false)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
}

func TestLogsShouldStreamAndFilterByContainer(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Container: testContainerID})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, false)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testContainerID, entry.ContainerId)
	}
}

func TestLogsShouldStreamAndFilterByNode(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Node: testNodeID})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, false)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Equal(t, testNodeID, entry.NodeId)
	}
}

func TestLogsShouldStreamAndFilterByService(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Service: "test"})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, false)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceID))
	}
}

func TestLogsShouldStreamAndFilterByMessage(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Message: testMessage})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, false)
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
	go produceLogEntries(100, false)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Message), testMessage)
	}
}

func TestLogsShouldStreamAndFilterByStack(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{Stack: "test"})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, false)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackName) || strings.HasPrefix(entry.StackId, testStackID))
	}
}

func TestLogsShouldStreamAndFilterByInfrastructureRole(t *testing.T) {
	stream, err := logsClient.GetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, false)
	entries, err := listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	for entry := range entries {
		assert.True(t, entry.Role != amp.InfrastructureRole)
	}
	stream.CloseSend()

	stream, err = logsClient.GetStream(ctx, &GetRequest{Infra: true})
	if err != nil {
		t.Error(err)
	}
	go produceLogEntries(100, true)
	entries, err = listenToLogEntries(stream, defaultNumberOfEntries)
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfEntries, len(entries))
	gotInfraEntry := false
	for entry := range entries {
		if entry.Role == amp.InfrastructureRole {
			gotInfraEntry = true
			break
		}
	}
	assert.True(t, gotInfraEntry)
}

func TestLogsEnd(t *testing.T) {
	sc.Close()
}

func produceLogEntries(howMany int, generateInfrasRole bool) error {
	for i := 0; i < howMany; i++ {
		logEntry := &LogEntry{
			ContainerId: testContainerID,
			Message:     testMessage + strconv.Itoa(rand.Int()),
			NodeId:      testNodeID,
			ServiceId:   testServiceID,
			ServiceName: testServiceName,
			StackId:     testStackID,
			StackName:   testStackName,
			TaskName:    testTaskName,
			TaskId:      testTaskID,
			Timestamp:   time.Now().Format(time.RFC3339Nano),
			TimeId:      time.Now().UTC().Format(time.RFC3339Nano),
		}
		if i%5 == 0 && generateInfrasRole {
			logEntry.Role = amp.InfrastructureRole
		}
		message, err := proto.Marshal(logEntry)
		if err != nil {
			return err
		}
		if err = sc.Publish(amp.LogsQueue, message); err != nil {
			return err
		}
	}
	return nil
}

func listenToLogEntries(stream Logs_GetStreamClient, howMany int) (chan *LogEntry, error) {
	entries := make(chan *LogEntry, howMany)
	entryCount := 0
	timeout := time.After(30 * time.Second)

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
