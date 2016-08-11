package server

import (
	"encoding/json"
	"github.com/appcelerator/amp/api/rpc/logs"
	"golang.org/x/net/context"
	"log"
)

const (
	esIndex = "amp-logs"
)

// logService is used to implement log.LogServer
type Logs struct {
}

// Get implements log.LogServer
func (s *Logs) Get(ctx context.Context, in *logs.GetRequest) (*logs.GetReply, error) {
	reply := logs.GetReply{}
	// Search with a term query
	//termQuery := elastic.NewTermQuery("user", "bquenin")
	searchResult, err := ES.GetClient().Search().
		Index(esIndex).
		//Query(termQuery).
		Sort("time_id", false).
		Size(100).
		Do()
	if err != nil {
		return &reply, err
	}
	reply.Entries = make([]*logs.LogEntry, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		var entry logs.LogEntry
		log.Printf("hit: %s", hit.Source)
		err := json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			log.Printf("h4")
			return &reply, err
		}
		reply.Entries[i] = &entry
	}
	log.Printf("reply: %+v\n", reply)
	return &reply, nil
}
