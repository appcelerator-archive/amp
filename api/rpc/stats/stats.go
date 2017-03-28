package stats

import (
	"fmt"
	"log"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// Stats structure to implement StatsServer interface
type Stats struct {
	ElasticsearchURL string
	EsConnected      bool
	Es               *elasticsearch.Elasticsearch
	Store            storage.Interface
	NatsStreaming    ns.NatsStreaming
	Docker           *dockerClient.Client
}

const (
	discriminatorContainer = "container"
	discriminatorService   = "service"
	discriminatorNode      = "node"
	discriminatorTask      = "task"
	metricsCPU             = "cpu"
	metricsMem             = "mem"
	metricsNet             = "net"
	metricsIO              = "io"
)

func (s *Stats) isElasticsearch(ctx context.Context) *elastic.Client {
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

func (s *Stats) doesElasticsearchServiceExist(ctx context.Context) bool {
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

// StatsQuery extracts stat information according to StatsRequest
func (s *Stats) StatsQuery(ctx context.Context, req *StatsRequest) (*StatsReply, error) {
	client := s.isElasticsearch(ctx)
	if client == nil {
		return nil, fmt.Errorf("the monitoring_elasticsearch service is not running, please start stack 'monitoring'")
	}

	// Prepare request to elasticsearch
	//request := client.Search().Index(EsIndex)
	return &StatsReply{}, nil
}
