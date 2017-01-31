package logs

import (
	"encoding/json"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"log"

	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/mq"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v3"
)

const (
	// EsIndex is the name of the Elasticsearch index
	EsIndex = "amp-logs"

	// EsType is the name of the Elasticsearch type
	EsType = "amp-log-entry"

	// EsMapping is the name of the Elasticsearch type mapping
	EsMapping = `{
	 "amp-log-entry": {
		"properties": {
		  "timestamp": {
			"type": "date"
		  },
		  "time_id": {
			"type": "keyword",
		  },
		  "container_id": {
			"type": "keyword",
		  },
		  "node_id": {
			"type": "keyword",
		  },
		  "service_id": {
			"type": "keyword",
		  },
		  "service_name": {
			"type": "keyword",
		  },
		  "task_id": {
			"type": "keyword",
		  },
		  "task_name": {
			"type": "keyword",
		  },
		  "stack_id": {
			"type": "keyword",
		  },
		  "stack_name": {
			"type": "keyword",
		  },
		  "role": {
			"type": "keyword",
		  }
		}
	  }
	}`
)

// Server is used to implement log.LogServer
type Server struct {
	Es    *elasticsearch.Elasticsearch
	Store storage.Interface
	MQ    mq.Interface
}

// Get implements log.LogServer
func (s *Server) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	// TODO: Authentication is disabled in order to allow tests. Re-enable this as soon as we have a way to auth in tests.
	//_, err := oauth.CheckAuthorization(ctx, logs.Store)
	//if err != nil {
	//	return nil, err
	//}

	log.Println("rpc-logs: Get", in.String())

	// Prepare request to elasticsearch
	request := s.Es.GetClient().Search().Index(EsIndex)
	request.Sort("time_id", false)
	if in.Size != 0 {
		request.Size(int(in.Size))
	} else {
		request.Size(100)
	}

	masterQuery := elastic.NewBoolQuery()
	if in.Container != "" {
		masterQuery.Filter(elastic.NewPrefixQuery("container_id", in.Container))
	}
	if in.Message != "" {
		queryString := elastic.NewSimpleQueryStringQuery(in.Message)
		queryString.Field("message")
		masterQuery.Filter(queryString)
	}
	if in.Node != "" {
		masterQuery.Filter(elastic.NewPrefixQuery("node_id", in.Node))
	}
	if in.Service != "" {
		boolQuery := elastic.NewBoolQuery()
		masterQuery.Filter(
			boolQuery.Should(elastic.NewPrefixQuery("service_id", in.Service)),
			boolQuery.Should(elastic.NewPrefixQuery("service_name", in.Service)),
		)
	}
	if in.Stack != "" {
		boolQuery := elastic.NewBoolQuery()
		masterQuery.Filter(
			boolQuery.Should(elastic.NewPrefixQuery("stack_id", in.Stack)),
			boolQuery.Should(elastic.NewPrefixQuery("stack_name", in.Stack)),
		)
	}
	if !in.Infra {
		masterQuery.MustNot(elastic.NewTermQuery("role", amp.InfrastructureRole))
	}
	// TODO timestamp queries

	// Perform request
	searchResult, err := request.Query(masterQuery).Do(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "%v", err)
	}

	// Build reply (from elasticsearch response)
	reply := GetReply{}
	reply.Entries = make([]*LogEntry, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		entry, err := parseJSONLogEntry(*hit.Source)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "%v", err)
		}
		reply.Entries[i] = &entry
	}

	// Reverse entries
	for i, j := 0, len(reply.Entries)-1; i < j; i, j = i+1, j-1 {
		reply.Entries[i], reply.Entries[j] = reply.Entries[j], reply.Entries[i]
	}
	log.Printf("rpc-logs: Get successful, returned %d entries\n", len(reply.Entries))
	return &reply, nil
}

// GetStream implements log.LogServer
func (s *Server) GetStream(in *GetRequest, stream Logs_GetStreamServer) error {
	log.Println("rpc-logs: GetStream", in.String())

	sub, err := s.MQ.Subscribe(amp.LogsQueue, func(msg proto.Message, err error) {
		if err != nil {
			log.Println("Error in message processing:", err)
			return
		}

		entry, ok := msg.(*LogEntry)
		if !ok {
			log.Println("Error in type assertion")
			return
		}
		if filter(entry, in) {
			stream.Send(entry)
		}
	}, &LogEntry{})
	if err != nil {
		sub.Unsubscribe()
		return grpc.Errorf(codes.Internal, "%v", err)
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
	if !in.Infra {
		match = entry.Role != amp.InfrastructureRole
	}
	return match
}
