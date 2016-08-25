package logs

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/appcelerator/amp/api/rpc/oauth"
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/kafka"
	"github.com/appcelerator/amp/data/storage"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v3"
	"strings"
)

const (
	esIndex       = "amp-logs"
	kafkaLogTopic = "amp-logs"
)

// Logs is used to implement log.LogServer
type Logs struct {
	Es    elasticsearch.Elasticsearch
	Store storage.Interface
	Kafka kafka.Kafka
}

// Get implements log.LogServer
func (logs *Logs) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	_, err := oauth.CheckAuthorization(ctx, logs.Store)
	if err != nil {
		return nil, err
	}
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
		request.Query(elastic.NewFuzzyQuery("message", in.Message))
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
	consumer, err := logs.Kafka.NewConsumer()
	if err != nil {
		return err
	}
	partitionConsumer, err := consumer.ConsumePartition(kafkaLogTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			entry, err := parseLogEntry(msg.Value)
			if err != nil {
				return err
			}
			if filter(&entry, in) {
				stream.Send(&entry)
			}

		case <-stream.Context().Done():
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
		match = entry.ServiceId == in.ServiceId
	}
	if in.ServiceName != "" {
		match = entry.ServiceName == in.ServiceName
	}
	if in.ContainerId != "" {
		match = entry.ContainerId == in.ContainerId
	}
	if in.NodeId != "" {
		match = entry.NodeId == in.NodeId
	}
	if in.Message != "" {
		match = strings.Contains(entry.Message, in.Message)
	}
	return match
}
