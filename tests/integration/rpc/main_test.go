package tests

import (
	"log"
	"os"
	"testing"
	//"time"

	"golang.org/x/net/context"
	//"google.golang.org/grpc"

	//tested packages
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"google.golang.org/grpc"
	"time"
)

var (
	ctx            context.Context
	functionClient function.FunctionClient
	statsClient    stats.StatsClient
	stackClient    stack.StackServiceClient
	topicClient    topic.TopicClient
	serviceClient  service.ServiceClient
	logsClient     logs.LogsClient
	store          storage.Interface
)

func TestMain(m *testing.M) {
	// Connect to amplifier
	log.Println("Connecting to amplifier")
	conn, err := grpc.Dial(amp.AmplifierDefaultEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second))
	if err != nil {
		log.Panicf("Unable to connect to amplifier on: %s\n%v", amp.AmplifierDefaultEndpoint, err)
	}
	log.Println("Connected to amplifier")

	store = etcd.New([]string{amp.EtcdDefaultEndpoint}, "amp")
	if err := store.Connect(5 * time.Second); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", amp.EtcdDefaultEndpoint, err)
	}

	ctx = context.Background()

	// init package clients
	functionClient = function.NewFunctionClient(conn)
	statsClient = stats.NewStatsClient(conn)
	stackClient = stack.NewStackServiceClient(conn)
	topicClient = topic.NewTopicClient(conn)
	serviceClient = service.NewServiceClient(conn)
	logsClient = logs.NewLogsClient(conn)

	// start tests
	os.Exit(m.Run())
}
