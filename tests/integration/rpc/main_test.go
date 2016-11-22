package tests

import (
	"log"
	"os"
	"testing"
	//"time"

	"github.com/appcelerator/amp/api/server"
	"golang.org/x/net/context"
	//"google.golang.org/grpc"

	//tested packages
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"google.golang.org/grpc"
	"time"
)

var (
	config        server.Config
	ctx           context.Context
	statsClient   stats.StatsClient
	stackClient   stack.StackServiceClient
	topicClient   topic.TopicClient
	serviceClient service.ServiceClient
	logsClient    logs.LogsClient
	store         storage.Interface
)

func TestMain(m *testing.M) {
	// Get configuration
	config = server.ConfigFromEnv()

	// Connect to amplifier
	log.Println("Connecting to amplifier")
	conn, err := grpc.Dial("amplifier:50101",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second))
	if err != nil {
		log.Panicln("Cannot connect to amplifier", err)
	}
	log.Println("Connected to amplifier")

	store = etcd.New(config.EtcdEndpoints, "amp")
	if err := store.Connect(5 * time.Second); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", config.EtcdEndpoints, err)
	}

	ctx = context.Background()

	// init package clients
	statsClient = stats.NewStatsClient(conn)
	stackClient = stack.NewStackServiceClient(conn)
	topicClient = topic.NewTopicClient(conn)
	serviceClient = service.NewServiceClient(conn)
	logsClient = logs.NewLogsClient(conn)

	// start tests
	os.Exit(m.Run())
}
