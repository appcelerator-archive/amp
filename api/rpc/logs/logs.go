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
	// NatsLogTopic is the topic used for logs
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
	// Prepare request to elasticsearch
	request := logs.Es.GetClient().Search().Index(esIndex)
	if in.From >= 0 {
		request.From(int(in.From))
	}
	if in.Size != 0 {
		request.Size(int(in.Size))
	} else {
		request.Size(100)
	}
	if in.ServiceId != "" {
		request.Query(elastic.NewTermQuery("service_id", in.ServiceId))
	}
	if in.ServiceName != "" {
		request.Query(elastic.NewTermQuery("service_name", in.ServiceName))
	}
	if in.ContainerId != "" {
		request.Query(elastic.NewTermQuery("container_id", in.ContainerId))
	}
	if in.NodeId != "" {
		request.Query(elastic.NewTermQuery("node_id", in.NodeId))
	}
	if in.Message != "" {
		queryString := elastic.NewQueryStringQuery("*" + in.Message)
		queryString.DefaultField("message")
		queryString.AnalyzeWildcard(true)
		request.Query(queryString)
	}
	// TODO timestamp queries

	// Perform request
	searchResult, err := request.Do()
	if err != nil {
		return nil, err
	}

	// Build reply (from elasticsearch response)
	reply := GetReply{}
	reply.Entries = make([]*LogEntry, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		entry, err := parseLogEntry(*hit.Source)
		if err != nil {
			return nil, err
		}
		reply.Entries[i] = &entry
	}
	return &reply, nil
}

// GetStream implements log.LogServer
func (logs *Logs) GetStream(in *GetRequest, stream Logs_GetStreamServer) error {
	sub, err := logs.Nats.Subscribe(NatsLogTopic, func(msg *stan.Msg) {
		logEntry := LogEntry{}
		err := proto.Unmarshal(msg.Data, &logEntry)
		if err != nil {
			log.Printf("error unmarshalling log entry: %v", err)
		}
		if filter(&logEntry, in) {
			stream.Send(&logEntry)
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

func parseLogEntry(data []byte) (LogEntry, error) {
	var entry LogEntry
	err := json.Unmarshal(data, &entry)
	if err != nil {
		return entry, err
	}
	return entry, err
}

func filter(entry *LogEntry, in *GetRequest) bool {
	match := true
	if in.ServiceId != "" {
		match = strings.EqualFold(entry.ServiceId, in.ServiceId)
	}
	if in.ServiceName != "" {
		match = strings.EqualFold(entry.ServiceName, in.ServiceName)
	}
	if in.ContainerId != "" {
		match = strings.EqualFold(entry.ContainerId, in.ContainerId)
	}
	if in.NodeId != "" {
		match = strings.EqualFold(entry.NodeId, in.NodeId)
	}
	if in.Message != "" {
		match = strings.Contains(strings.ToLower(entry.Message), strings.ToLower(in.Message))
	}
	return match
}
