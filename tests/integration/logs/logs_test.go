package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/pkg/labels"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/appcelerator/amp/tests"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	testMessage            = "test message "
	testContainerID        = stringid.GenerateNonCryptoID()
	testContainerName      = "testcontainer"
	testContainerShortName = "testcontainershortname"
	testContainerState     = "testcontainerstate"
	testServiceID          = stringid.GenerateNonCryptoID()
	testServiceName        = "testservice"
	testStackName          = "teststack"
	testNodeID             = stringid.GenerateNonCryptoID()
	testTaskID             = stringid.GenerateNonCryptoID()
)

var (
	ctx         context.Context
	client      LogsClient
	credentials metadata.MD
	lp          *LogProducer
)

func setup() (err error) {
	if credentials, err = helpers.Login(); err != nil {
		return err
	}
	conn, err := helpers.AmplifierConnection()
	if err != nil {
		return err
	}
	client = NewLogsClient(conn)
	lp = NewLogProducer()
	ctx = metadata.NewContext(context.Background(), credentials)

	// Populate logs
	if err := lp.produce(NumberOfEntries); err != nil {
		return err
	}
	for {
		time.Sleep(1 * time.Second)
		r, err := client.Get(ctx, &GetRequest{Service: testServiceID})
		if err != nil {
			log.Println(err)
			continue
		}
		if len(r.Entries) == NumberOfEntries {
			break
		}
	}
	return nil
}

func tearDown() {
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatalln(err)
	}
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestLogsShouldGetAHundredLogEntriesByDefault(t *testing.T) {
	expected := NumberOfEntries
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

func TestLogsShouldFilterByContainer(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Container: testContainerID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, testContainerID, entry.ContainerId)
	}
}

func TestLogsShouldFilterByNode(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Node: testNodeID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, testNodeID, entry.NodeId)
	}
}

func TestLogsShouldFilterByService(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Service: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceID))
	}

	r, err = client.Get(ctx, &GetRequest{Service: testServiceName})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceID))
	}
}

func TestLogsShouldFilterByMessage(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Message: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Contains(t, strings.ToLower(entry.Msg), "test")
	}
}

func TestLogsShouldFilterByStackName(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Stack: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackName))
	}
}

func TestLogsShouldExcludeAmpLogs(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Empty(t, entry.Labels[labels.KeyRole])
	}
}

func TestLogsShouldIncludeAmpLogs(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Service: "amp_amplifier", IncludeAmpLogs: true})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	gotInfraEntry := false
	for _, entry := range r.Entries {
		if entry.Labels[labels.KeyRole] != "" {
			gotInfraEntry = true
			break
		}
	}
	assert.True(t, gotInfraEntry)
}

func TestLogsShouldFetchGivenNumberOfEntries(t *testing.T) {
	for i := int64(1); i < 100; i += 10 {
		r, err := client.Get(ctx, &GetRequest{Size: i})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, i, int64(len(r.Entries)))
	}
}

func TestLogsShouldBeOrdered(t *testing.T) {
	r, err := client.Get(ctx, &GetRequest{Container: testContainerID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	var current, previous int64
	for _, entry := range r.Entries {
		current, err = strconv.ParseInt(strings.TrimPrefix(entry.Msg, testMessage), 16, 64)
		assert.NoError(t, err)
		assert.True(t, current > previous, "Should be true but got current: %v <= previous: %v", current, previous)
		previous = current
	}
}

func TestLogsShouldStreamLogs(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))
}

func TestLogsShouldStreamAndFilterByContainer(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{Container: testContainerID})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Equal(t, testContainerID, entry.ContainerId)
	}
}

func TestLogsShouldStreamAndFilterByNode(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{Node: testNodeID})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Equal(t, testNodeID, entry.NodeId)
	}
}

func TestLogsShouldStreamAndFilterByService(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{Service: "test"})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, testServiceName) || strings.HasPrefix(entry.ServiceId, testServiceID))
	}
}

func TestLogsShouldStreamAndFilterByMessage(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{Message: testMessage})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Msg), testMessage)
	}
}

func TestLogsShouldStreamAndFilterCaseInsensitivelyByMessage(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{Message: strings.ToUpper(testMessage)})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Msg), testMessage)
	}
}

func TestLogsShouldStreamAndFilterByStackName(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{Stack: "test"})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.StackName, testStackName))
	}
}

func TestLogsShouldStreamAndExcludeAmpLogs(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.True(t, entry.Labels[labels.KeyRole] != labels.ValueRoleInfrastructure)
	}
}

func TestLogsShouldStreamAndIncludeAmpLogs(t *testing.T) {
	conn, err := helpers.AmplifierConnection()
	assert.NoError(t, err)
	defer conn.Close()
	client = NewLogsClient(conn)

	stream, err := client.GetStream(ctx, &GetRequest{IncludeAmpLogs: true})
	if err != nil {
		t.Error(err)
	}

	lp.startAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.stopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	gotInfraEntry := false
	for entry := range entries {
		if entry.Labels[labels.KeyRole] == labels.ValueRoleInfrastructure {
			gotInfraEntry = true
			break
		}
	}
	assert.True(t, gotInfraEntry)
}

// Helpers

func listenToLogEntries(stream Logs_GetStreamClient, howMany int) (chan *LogEntry, error) {
	entries := make(chan *LogEntry, howMany)
	entryCount := 0
	timeout := time.After(30 * time.Second)

	defer close(entries)

	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		select {
		case entries <- entry:
			entryCount++
			if entryCount == howMany {
				return entries, nil
			}
		case <-timeout:
			return entries, nil
		}
	}
}

type LogProducer struct {
	ns              *ns.NatsStreaming
	asyncProduction int32
	counter         int64
}

func NewLogProducer() *LogProducer {
	lp := &LogProducer{
		ns:      ns.NewClient(ns.DefaultURL, ns.ClusterID, stringid.GenerateNonCryptoID(), 60*time.Second),
		counter: 0,
	}
	if err := lp.ns.Connect(); err != nil {
		log.Fatalln("Cannot connect to NATS", err)
	}
	go func(lp *LogProducer) {
		for {
			time.Sleep(50 * time.Millisecond)
			if lp.asyncProduction > 0 {
				if err := lp.produce(NumberOfEntries); err != nil {
					log.Println("error producing async messages", err)
				}
			}
		}
	}(lp)
	return lp
}

func (lp *LogProducer) buildLogEntry(infrastructure bool) *LogEntry {
	atomic.AddInt64(&lp.counter, 1)
	entry := &LogEntry{
		Timestamp:          time.Now().UTC().Format(time.RFC3339Nano),
		ContainerId:        testContainerID,
		ContainerName:      testContainerName,
		ContainerShortName: testContainerShortName,
		ContainerState:     testContainerState,
		ServiceName:        testServiceName,
		ServiceId:          testServiceID,
		TaskId:             testTaskID,
		StackName:          testStackName,
		NodeId:             testNodeID,
		TimeId:             fmt.Sprintf("%016X", lp.counter),
		Labels:             make(map[string]string),
		Msg:                testMessage + fmt.Sprintf("%016X", lp.counter),
	}
	if infrastructure {
		entry.Labels[labels.KeyRole] = labels.ValueRoleInfrastructure
	}
	return entry
}

func (lp *LogProducer) produce(howMany int) error {
	entries := GetReply{}
	for i := 0; i < howMany; i++ {
		// User log entry
		user := lp.buildLogEntry(false)
		entries.Entries = append(entries.Entries, user)

		// Infrastructure log entry
		infra := lp.buildLogEntry(true)
		entries.Entries = append(entries.Entries, infra)
	}
	message, err := proto.Marshal(&entries)
	if err != nil {
		return err
	}
	if err := lp.ns.GetClient().Publish(ns.LogsSubject, message); err != nil {
		return err
	}
	return nil
}

func (lp *LogProducer) startAsyncProducer() {
	atomic.CompareAndSwapInt32(&lp.asyncProduction, 0, 1)
}

func (lp *LogProducer) stopAsyncProducer() {
	atomic.CompareAndSwapInt32(&lp.asyncProduction, 1, 0)
}
