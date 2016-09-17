package stack_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	defaultPort             = ":50101"
	etcdDefaultEndpoints    = "http://localhost:2379"
	serverAddress           = "localhost" + defaultPort
	elasticsearchDefaultURL = "http://localhost:9200"
	kafkaDefaultURL         = "localhost:9092"
	influxDefaultURL        = "http://localhost:8086"
	example                 = `web:
  image: appcelerator.io/amp-demo
  public:
    - name: www
      protocol: tcp
      publish_port: 90
      internal_port: 3000
  replicas: 3
  environment:
    REDIS_PASSWORD: password
redis:
  image: redis
  environment:
    - PASSWORD=password`
)

var (
	config           server.Config
	port             string
	etcdEndpoints    string
	elasticsearchURL string
	kafkaURL         string
	influxURL        string
	client           stack.StackServiceClient
	ctx              context.Context
)

func parseEnv() {
	port = os.Getenv("port")
	if port == "" {
		port = defaultPort
	}
	etcdEndpoints = os.Getenv("endpoints")
	if etcdEndpoints == "" {
		etcdEndpoints = etcdDefaultEndpoints
	}
	elasticsearchURL = os.Getenv("elasticsearchURL")
	if elasticsearchURL == "" {
		elasticsearchURL = elasticsearchDefaultURL
	}
	kafkaURL = os.Getenv("kafkaURL")
	if kafkaURL == "" {
		kafkaURL = kafkaDefaultURL
	}
	influxURL = os.Getenv("influxURL")
	if influxURL == "" {
		influxURL = influxDefaultURL
	}

	// update config
	config.Port = port
	for _, s := range strings.Split(etcdEndpoints, ",") {
		config.EtcdEndpoints = append(config.EtcdEndpoints, s)
	}
	config.ElasticsearchURL = elasticsearchURL
	config.KafkaURL = kafkaURL
	config.InfluxURL = influxURL
}

func TestMain(m *testing.M) {
	parseEnv()
	go server.Start(config)

	ctx = context.Background()

	// there is no event when the server starts listening, so we just wait a second
	time.Sleep(1 * time.Second)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		fmt.Println("connection failure")
		os.Exit(1)
	}
	client = stack.NewStackServiceClient(conn)
	os.Exit(m.Run())
}

func TestShouldUpStackSuccessfully(t *testing.T) {
	r, err := client.Up(ctx, &stack.UpRequest{Stackfile: example})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, r.StackId, "StackId should not be empty")
}
