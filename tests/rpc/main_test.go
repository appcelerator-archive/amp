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
)

var (
	config        server.Config
	ctx           context.Context
	statsClient   stats.StatsClient
	stackClient   stack.StackServiceClient
	topicClient   topic.TopicClient
	serviceClient service.ServiceClient
	logsClient    logs.LogsClient
)

func TestMain(m *testing.M) {

	//init server and context
	conf, conn := server.StartTestServer()
	config = conf
	//config = server.GetConfig()
	/*
		log.Printf("Connecting to amplifier %s:%s\n", config.ServerAddress, config.ServerPort)
		conn, err := grpc.Dial(config.ServerAddress+":"+config.ServerPort,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(60*time.Second))
		if err != nil {
			log.Panicln("Cannot connect to amplifier", err)
			os.Exit(1)
		}
	*/
	log.Println("Connected to amplifier")
	ctx = context.Background()

	//init package clients
	statsClient = stats.NewStatsClient(conn)
	stackClient = stack.NewStackServiceClient(conn)
	topicClient = topic.NewTopicClient(conn)
	serviceClient = service.NewServiceClient(conn)
	logsClient = logs.NewLogsClient(conn)

	//start tests
	os.Exit(m.Run())
}
