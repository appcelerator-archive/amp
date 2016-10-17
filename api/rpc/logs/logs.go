package logs

import (
	"encoding/json"
	"strings"

	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v3"
	"log"
)

const (
	esIndex = "amp-logs"
	// NatsLogTopic is the nats channel
	NatsLogTopic = "amp-logs"
)

// Logs is used to implement log.LogServer
type Logs struct {
	Es    elasticsearch.Elasticsearch
	Store storage.Interface
	Nats  stan.Conn
}

// Get implements log.LogServer
func (logs *Logs) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	// TODO: Authentication is disabled in order to allow tests. Re-enable this as soon as we have a way to auth in tests.
	//_, err := oauth.CheckAuthorization(ctx, logs.Store)
	//if err != nil {
	//	return nil, err
	//}
	log.Printf("log requested: [%v]", in)
	// Prepare request to elasticsearch
	request := logs.Es.GetClient().Search().Index(esIndex)
	request.Sort("time_id", false)
	if in.Size != 0 {
		request.Size(int(in.Size))
	} else {
		request.Size(100)
	}

	masterQuery := elastic.NewBoolQuery()
	if in.Container != "" {
		masterQuery.Must(elastic.NewPrefixQuery("container_id", in.Container))
	}
	if in.Message != "" {
		queryString := elastic.NewQueryStringQuery(in.Message + "*")
		queryString.Field("message")
		queryString.AnalyzeWildcard(true)
		masterQuery.Must(queryString)
	}
	if in.Node != "" {
		masterQuery.Must(elastic.NewPrefixQuery("node_id", in.Node))
	}
	if in.Service != "" {
		masterQuery.Should(elastic.NewPrefixQuery("service_id", in.Service))
		masterQuery.Should(elastic.NewPrefixQuery("service_name", in.Service))
	}
	if in.Stack != "" {
		masterQuery.Should(elastic.NewPrefixQuery("stack_id", in.Stack))
		masterQuery.Should(elastic.NewPrefixQuery("stack_name", in.Stack))
	}
	// TODO timestamp queries

	// Perform request
	searchResult, err := request.Query(masterQuery).Do()
	if err != nil {
		return nil, err
	}

	// Build reply (from elasticsearch response)
	reply := GetReply{}
	reply.Entries = make([]*LogEntry, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		entry, err := parseJSONLogEntry(*hit.Source)
		if err != nil {
			return nil, err
		}
		reply.Entries[i] = &entry
	}

	// Reverse entries
	for i, j := 0, len(reply.Entries)-1; i < j; i, j = i+1, j-1 {
		reply.Entries[i], reply.Entries[j] = reply.Entries[j], reply.Entries[i]
	}

	return &reply, nil
}

// GetStream implements log.LogServer
func (logs *Logs) GetStream(in *GetRequest, stream Logs_GetStreamServer) error {
	log.Printf("log stream requested: [%v]", in)
	sub, err := logs.Nats.Subscribe(NatsLogTopic, func(msg *stan.Msg) {
		entry, err := parseProtoLogEntry(msg.Data)
		if err != nil {
			return
		}
		if filter(&entry, in) {
			stream.Send(&entry)
		}
	})
	if err != nil {
		sub.Unsubscribe()
		return err
	}
	for {
		select {
		case <-stream.Context().Done():
			sub.Unsubscribe()
			return stream.Context().Err()
		}
	}
}

func parseJSONLogEntry(data []byte) (logEntry LogEntry, err error) {
	err = json.Unmarshal(data, &logEntry)
	return
}

func parseProtoLogEntry(data []byte) (logEntry LogEntry, err error) {
	err = proto.Unmarshal(data, &logEntry)
	return
}

func filter(entry *LogEntry, in *GetRequest) bool {
	match := true
	if in.Container != "" {
		match = strings.HasPrefix(strings.ToLower(entry.ContainerId), strings.ToLower(in.Container))
	}
	if in.Message != "" {
		match = strings.Contains(strings.ToLower(entry.Message), strings.ToLower(in.Message))
	}
	if in.Node != "" {
		match = strings.HasPrefix(strings.ToLower(entry.NodeId), strings.ToLower(in.Node))
	}
	if in.Service != "" {
		serviceID := strings.ToLower(entry.ServiceId)
		serviceName := strings.ToLower(entry.ServiceName)
		match = strings.HasPrefix(serviceID, strings.ToLower(in.Service)) || strings.HasPrefix(serviceName, strings.ToLower(in.Service))
	}
	if in.Stack != "" {
		stackID := strings.ToLower(entry.StackId)
		stackName := strings.ToLower(entry.StackName)
		match = strings.HasPrefix(stackID, strings.ToLower(in.Stack)) || strings.HasPrefix(stackName, strings.ToLower(in.Stack))
	}
	return match
}
