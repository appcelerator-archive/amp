package server

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/logs"
	"golang.org/x/net/context"
	"reflect"
)

const (
	esIndex = "amp-logs"
)

// logService is used to implement log.LogServer
type logsService struct {
}

// Get implements log.LogServer
func (s *logsService) Get(ctx context.Context, in *logs.GetRequest) (*logs.GetReply, error) {
	// Search with a term query
	//termQuery := elastic.NewTermQuery("user", "bquenin")
	searchResult, err := es.GetClient().Search().
		Index(esIndex).
		//Query(termQuery).
		Sort("time_id", false).
		Size(100).
		Do()

	if err != nil {
		// Handle error
		panic(err)
	}
	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var entry logs.LogEntry
	for _, item := range searchResult.Each(reflect.TypeOf(entry)) {
		if e, ok := item.(logs.LogEntry); ok {
			fmt.Printf("Message: %s\n", e.Message)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d log entries\n", searchResult.TotalHits())

	//
	//// Here's how you iterate through results with full control over each step.
	//if searchResult.Hits.TotalHits > 0 {
	//	fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)
	//
	//	// Iterate through results
	//	for _, hit := range searchResult.Hits.Hits {
	//		// hit.Index contains the name of the index
	//
	//		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
	//		var t Tweet
	//		err := json.Unmarshal(*hit.Source, &t)
	//		if err != nil {
	//			// Deserialization failed
	//		}
	//
	//		// Work with tweet
	//		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	//	}
	//} else {
	//	// No hits
	//	fmt.Print("Found no tweets\n")
	//}
	return &logs.GetReply{}, nil
}
