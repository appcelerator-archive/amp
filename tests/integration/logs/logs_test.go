package logs

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/cluster/constants"
	. "github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	ctx context.Context
	h   *helpers.Helper
	lp  *helpers.LogProducer
)

func setup() (err error) {
	// Test helper
	if h, err = helpers.New(); err != nil {
		return err
	}

	// Login context
	credentials, err := h.Login()
	if err != nil {
		return err
	}
	ctx = metadata.NewOutgoingContext(context.Background(), credentials)

	// Log producer helper
	lp = helpers.NewLogProducer(h)
	if err := lp.PopulateLogs(); err != nil {
		return err
	}
	return nil
}

func tearDown() {
}

//client = NewLogsClient(conn)

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
		r, err := h.Logs().LogsGet(ctx, &GetRequest{})
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
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Container: helpers.TestContainerID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, helpers.TestContainerID, entry.ContainerId)
	}
}

func TestLogsShouldFilterByNode(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Node: helpers.TestNodeID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Equal(t, helpers.TestNodeID, entry.NodeId)
	}
}

func TestLogsShouldFilterByService(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Service: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, helpers.TestServiceName) || strings.HasPrefix(entry.ServiceId, helpers.TestServiceID))
	}

	r, err = h.Logs().LogsGet(ctx, &GetRequest{Service: helpers.TestServiceID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, helpers.TestServiceName) || strings.HasPrefix(entry.ServiceId, helpers.TestServiceID))
	}
}

func TestLogsShouldFilterByMessage(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Message: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Contains(t, strings.ToLower(entry.Msg), "test")
	}
}

func TestLogsShouldFilterByStack(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Stack: "test"})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, helpers.TestStackName) || strings.HasPrefix(entry.StackId, helpers.TestStackID))
	}

	r, err = h.Logs().LogsGet(ctx, &GetRequest{Stack: helpers.TestStackID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.True(t, strings.HasPrefix(entry.StackName, helpers.TestStackName) || strings.HasPrefix(entry.StackId, helpers.TestStackID))
	}
}

func TestLogsShouldExcludeAmpLogs(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	for _, entry := range r.Entries {
		assert.Empty(t, entry.Labels[constants.LabelKeyRole])
	}
}

func TestLogsShouldIncludeAmpLogs(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Service: "amp_amplifier", IncludeAmpLogs: true})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	gotInfraEntry := false
	for _, entry := range r.Entries {
		if entry.Labels[constants.LabelKeyRole] != "" {
			gotInfraEntry = true
			break
		}
	}
	assert.True(t, gotInfraEntry)
}

func TestLogsShouldFetchGivenNumberOfEntries(t *testing.T) {
	for i := int32(1); i < 100; i += 10 {
		r, err := h.Logs().LogsGet(ctx, &GetRequest{Size: i})
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, i, int32(len(r.Entries)))
	}
}

func TestLogsShouldBeOrdered(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Container: helpers.TestContainerID})
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, r.Entries, "We should have at least one entry")
	var current, previous int64
	for _, entry := range r.Entries {
		current, err = strconv.ParseInt(strings.TrimPrefix(entry.Msg, helpers.TestMessage), 16, 64)
		assert.NoError(t, err)
		assert.True(t, current > previous, "Should be true but got current: %v <= previous: %v", current, previous)
		previous = current
	}
}

func TestLogsShouldPaginate(t *testing.T) {
	r, err := h.Logs().LogsGet(ctx, &GetRequest{Container: helpers.TestContainerID, Size: 20})
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, r.Entries, 20, "We should have exactly 20 entries")

	// Keep track of the IDs
	IDs := []string{}
	for _, entry := range r.Entries {
		IDs = append(IDs, entry.FromId)
	}

	// Note that the log entries are reversed compared to the ES query in order to display logs from least recent to most recent.
	// Therefore, We need to use the first result FromId (the least recent one) as the starting point in subsequent request
	r, err = h.Logs().LogsGet(ctx, &GetRequest{Container: helpers.TestContainerID, Size: 10, From: IDs[0]})
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, r.Entries, 10, "We should have exactly 10 entries")

	// Make sure we get different entries
	for _, entry := range r.Entries {
		assert.NotContains(t, IDs, entry.FromId)
	}
}

func TestLogsShouldStreamLogs(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))
}

func TestLogsShouldStreamAndFilterByContainer(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{Container: helpers.TestContainerID})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Equal(t, helpers.TestContainerID, entry.ContainerId)
	}
}

func TestLogsShouldStreamAndFilterByNode(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{Node: helpers.TestNodeID})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Equal(t, helpers.TestNodeID, entry.NodeId)
	}
}

func TestLogsShouldStreamAndFilterByService(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{Service: "test"})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.ServiceName, helpers.TestServiceName) || strings.HasPrefix(entry.ServiceId, helpers.TestServiceID))
	}
}

func TestLogsShouldStreamAndFilterByMessage(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{Message: helpers.TestMessage})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Msg), helpers.TestMessage)
	}
}

func TestLogsShouldStreamAndFilterCaseInsensitivelyByMessage(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{Message: strings.ToUpper(helpers.TestMessage)})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.Contains(t, strings.ToLower(entry.Msg), helpers.TestMessage)
	}
}

func TestLogsShouldStreamAndFilterByStack(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{Stack: "test"})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.True(t, strings.HasPrefix(entry.StackName, helpers.TestStackName) || strings.HasPrefix(entry.StackId, helpers.TestStackID))
	}
}

func TestLogsShouldStreamAndExcludeAmpLogs(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	for entry := range entries {
		assert.NotContains(t, entry.Labels, constants.LabelKeyRole)
	}
}

func TestLogsShouldStreamAndIncludeAmpLogs(t *testing.T) {
	stream, err := h.Logs().LogsGetStream(ctx, &GetRequest{IncludeAmpLogs: true})
	if err != nil {
		t.Error(err)
	}

	lp.StartAsyncProducer()
	entries, err := listenToLogEntries(stream, NumberOfEntries)
	lp.StopAsyncProducer()
	assert.NoError(t, err)
	assert.Equal(t, NumberOfEntries, len(entries))

	gotInfraEntry := false
	for entry := range entries {
		if _, exists := entry.Labels[constants.LabelKeyRole]; exists {
			gotInfraEntry = true
			break
		}
	}
	assert.True(t, gotInfraEntry)
}

// Helpers

func listenToLogEntries(stream Logs_LogsGetStreamClient, howMany int) (chan *LogEntry, error) {
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
