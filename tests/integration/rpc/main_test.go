package tests

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/topic"
	as "github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"
)

var (
	ctx            context.Context
	store          storage.Interface
	functionClient function.FunctionClient
	statsClient    stats.StatsClient
	stackClient    stack.StackServiceClient
	topicClient    topic.TopicClient
	serviceClient  service.ServiceClient
	logsClient     logs.LogsClient
	accountClient  account.AccountClient
	accountStore   as.Interface
)

func TestMain(m *testing.M) {
	ctx = context.Background()

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

	// Stores
	store = etcd.New([]string{amp.EtcdDefaultEndpoint}, "amp")
	if err := store.Connect(5 * time.Second); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", amp.EtcdDefaultEndpoint, err)
	}
	accountStore = as.NewStore(store)

	// Clients init
	functionClient = function.NewFunctionClient(conn)
	statsClient = stats.NewStatsClient(conn)
	stackClient = stack.NewStackServiceClient(conn)
	topicClient = topic.NewTopicClient(conn)
	serviceClient = service.NewServiceClient(conn)
	logsClient = logs.NewLogsClient(conn)
	accountClient = account.NewAccountClient(conn)

	// Start tests
	os.Exit(m.Run())
}
