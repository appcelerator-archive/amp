package server

import (
	"encoding/json"
	"github.com/appcelerator/amp/api/rpc/logs"
	"golang.org/x/net/context"
)

const (
	esIndex = "amp-logs"
)

// logService is used to implement log.LogServer
type Logs struct {
}

// Get implements log.LogServer
func (s *Logs) Get(ctx context.Context, in *logs.GetRequest) (*logs.GetReply, error) {
	// Search with a term query
	//termQuery := elastic.NewTermQuery("user", "bquenin")
	searchResult, err := ES.GetClient().Search().
		Index(esIndex).
		//Query(termQuery).
		Sort("time_id", false).
		Size(100).
		Do()
	if err != nil {
		return nil, err
	}

	reply := logs.GetReply{}
	reply.Entries = make([]*logs.LogEntry, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		var entry logs.LogEntry
		err := json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return nil, err
		}
		reply.Entries[i] = &entry
	}
	return &reply, nil
}
