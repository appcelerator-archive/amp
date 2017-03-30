package logs

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	ElasticsearchURL string
	EsConnected      bool
	Es               *elasticsearch.Elasticsearch
	Store            storage.Interface
	NatsStreaming    ns.NatsStreaming
	Docker           *dockerClient.Client
}

func (s *Server) isElasticsearch(ctx context.Context) *elastic.Client {
	if !s.doesElasticsearchServiceExist(ctx) {
		s.EsConnected = false
		return nil
	}
	if !s.EsConnected {
		log.Println("Connecting to elasticsearch at", s.ElasticsearchURL)
		if err := s.Es.Connect(s.ElasticsearchURL, amp.DefaultTimeout); err != nil {
			log.Printf("unable to connect to elasticsearch at %s: %v", s.ElasticsearchURL, err)
			return nil
		}
		s.EsConnected = true
		log.Println("Connected to elasticsearch at", s.ElasticsearchURL)
	}
	client := s.Es.GetClient()
	if client.IsRunning() {
		return client
	}
	return nil
}

func (s *Server) doesElasticsearchServiceExist(ctx context.Context) bool {
	list, err := s.Docker.ServiceList(ctx, types.ServiceListOptions{
	//Filter: filter,
	})
	if err != nil || len(list) == 0 {
		return false
	}
	for _, serv := range list {
		if serv.Spec.Annotations.Name == "monitoring_elasticsearch" {
			return true
		}
	}
	return false
}

// Get implements log.LogServer
func (s *Server) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	// TODO: Authentication is disabled in order to allow tests. Re-enable this as soon as we have a way to auth in tests.
	//_, err := oauth.CheckAuthorization(ctx, logs.Store)
	//if err != nil {
	//	return nil, err
	//}
	client := s.isElasticsearch(ctx)
	if client == nil {
		return nil, fmt.Errorf("the monitoring_elasticsearch service is not running, please start stack 'monitoring'")
	}
	log.Println("rpc-logs: Get", in.String())

	// Prepare request to elasticsearch
	request := client.Search().Index(EsIndex)
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

	sub, err := s.NatsStreaming.GetClient().Subscribe(amp.NatsLogsTopic, func(msg *stan.Msg) {
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
