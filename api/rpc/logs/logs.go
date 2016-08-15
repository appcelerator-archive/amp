package logs

import (
	"encoding/json"
	"github.com/appcelerator/amp/data/elasticsearch"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v3"
)

const (
	esIndex = "amp-logs"
)

// Logs is used to implement log.LogServer
type Logs struct {
	ES elasticsearch.Elasticsearch
}

// Get implements log.LogServer
func (s *Logs) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	// Prepare request to elasticsearch
	request := s.ES.GetClient().Search().Index(esIndex)
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
		var entry LogEntry
		err := json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return nil, err
		}
		reply.Entries[i] = &entry
	}
	return &reply, nil
}
