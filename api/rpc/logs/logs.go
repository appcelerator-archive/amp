package logs

import (
	"encoding/json"
	"github.com/appcelerator/amp/data/elasticsearch"
	"golang.org/x/net/context"
	"log"
)

const (
	esIndex = "amp-logs"
)

// logService is used to implement log.LogServer
type Logs struct {
	ES elasticsearch.Elasticsearch
}

// Get implements log.LogServer
func (s *Logs) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	reply := GetReply{}
	// Search with a term query
	//termQuery := elastic.NewTermQuery("user", "bquenin")
	searchResult, err := s.ES.GetClient().Search().
		Index(esIndex).
		//Query(termQuery).
		Sort("time_id", false).
		Size(100).
		Do()
	if err != nil {
		return &reply, err
	}
	reply.Entries = make([]*LogEntry, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		var entry LogEntry
		log.Printf("hit: %s", hit.Source)
		err := json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return &reply, err
		}
		reply.Entries[i] = &entry
	}
	return &reply, nil
}
