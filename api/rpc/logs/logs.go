package logs

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/elasticsearch"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/olivere/elastic.v5"
)

const (
	// InfrastructureRole amp role InfrastructureRole
	InfrastructureRole = "infrastructure"

	// ToolsRole amp role tools
	ToolsRole = "tools"
)

// Server is used to implement log.LogServer
type Server struct {
	Docker        *docker.Docker
	Es            *elasticsearch.Elasticsearch
	NatsStreaming *ns.NatsStreaming
}

// Get implements logs.LogsServer
func (s *Server) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	if err := s.Es.Connect(); err != nil {
		return nil, errors.New("unable to connect to elasticsearch service")
	}
	log.Println("rpc-logs: Get", in.String())

	// Prepares indices
	indices := []string{}
	date := time.Now().UTC()
	for i := 0; i < 2; i++ {
		indices = append(indices, "ampbeat-"+date.Format("2006.01.02"))
		date = date.AddDate(0, 0, -1)
	}

	// Prepare request to elasticsearch
	request := s.Es.GetClient().Search().Index(indices...).IgnoreUnavailable(true)
	request.Type("logs")
	request.Sort("@timestamp", false)
	if in.Size != 0 {
		request.Size(int(in.Size))
	} else {
		request.Size(100)
	}

	masterQuery := elastic.NewBoolQuery()
	if in.Container != "" {
		boolQuery := elastic.NewBoolQuery()
		masterQuery.Filter(
			boolQuery.Should(elastic.NewPrefixQuery("container_id", in.Service)),
			boolQuery.Should(elastic.NewPrefixQuery("container_name", in.Service)),
		)
	}
	if in.Service != "" {
		boolQuery := elastic.NewBoolQuery()
		masterQuery.Filter(
			boolQuery.Should(elastic.NewPrefixQuery("service_id", in.Service)),
			boolQuery.Should(elastic.NewPrefixQuery("service_name", in.Service)),
		)
	}
	if in.Stack != "" {
		masterQuery.Filter(elastic.NewPrefixQuery("stack_name", in.Stack))
	}
	if in.Node != "" {
		masterQuery.Filter(elastic.NewPrefixQuery("node_id", in.Node))
	}
	if in.Message != "" {
		queryString := elastic.NewSimpleQueryStringQuery(in.Message)
		queryString.Field("msg")
		masterQuery.Filter(queryString)
	}
	if !in.Infra {
		masterQuery.MustNot(elastic.NewTermQuery("role", InfrastructureRole))
		//For now ToolsRole is manage as InfrastructureRole
		//Later: ToolsRole should be manage with a user premission (admin)
		masterQuery.MustNot(elastic.NewTermQuery("role", ToolsRole))
	}

	// Perform request
	searchResult, err := request.Query(masterQuery).Do(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "%v", err)
	}

	// Build reply (from elasticsearch response)
	reply := GetReply{}
	reply.Entries = make([]*LogEntry, len(searchResult.Hits.Hits))
	for i, hit := range searchResult.Hits.Hits {
		entry := &LogEntry{}
		if err := json.Unmarshal(*hit.Source, &entry); err != nil {
			return nil, grpc.Errorf(codes.Internal, "%v", err)
		}
		reply.Entries[i] = entry
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
	if err := s.NatsStreaming.Connect(); err != nil {
		return errors.New("unable to connect to nats service")
	}
	log.Println("rpc-logs: GetStream", in.String())

	sub, err := s.NatsStreaming.GetClient().Subscribe(ns.LogsSubject, func(msg *stan.Msg) {
		entry := &LogEntry{}
		if err := proto.Unmarshal(msg.Data, entry); err != nil {
			return
		}
		if filter(entry, in) {
			stream.Send(entry)
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

func filter(entry *LogEntry, in *GetRequest) bool {
	match := true
	if in.Container != "" {
		containerID := strings.ToLower(entry.ContainerId)
		containerName := strings.ToLower(entry.ContainerName)
		match = strings.HasPrefix(containerID, strings.ToLower(in.Container)) || strings.HasPrefix(containerName, strings.ToLower(in.Container))
	}
	if in.Service != "" {
		serviceID := strings.ToLower(entry.ServiceId)
		serviceName := strings.ToLower(entry.ServiceName)
		match = strings.HasPrefix(serviceID, strings.ToLower(in.Service)) || strings.HasPrefix(serviceName, strings.ToLower(in.Service))
	}
	if in.Stack != "" {
		match = strings.HasPrefix(strings.ToLower(entry.StackName), strings.ToLower(in.Stack))
	}
	if in.Node != "" {
		match = strings.HasPrefix(strings.ToLower(entry.NodeId), strings.ToLower(in.Node))
	}
	if in.Message != "" {
		match = strings.Contains(strings.ToLower(entry.Msg), strings.ToLower(in.Message))
	}
	if !in.Infra {
		match = entry.Role != InfrastructureRole && entry.Role != ToolsRole
	}
	return match
}
