package tests

import (
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/topic"
	as "github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/account/schema"
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
	ctx                    context.Context
	store                  storage.Interface
	functionClient         function.FunctionClient
	statsClient            stats.StatsClient
	stackClient            stack.StackServiceClient
	topicClient            topic.TopicClient
	serviceClient          service.ServiceClient
	logsClient             logs.LogsClient
	accountClient          account.AccountClient
	anonymousAccountClient account.AccountClient
	accountStore           as.Interface
	testUser               = schema.User{
		Email:        "test@amp.io",
		Name:         "testUser",
		IsVerified:   true,
		PasswordHash: "testHash",
	}
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Stores
	store = etcd.New([]string{amp.EtcdDefaultEndpoint}, "amp")
	if err := store.Connect(5 * time.Second); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", amp.EtcdDefaultEndpoint, err)
	}
	accountStore = as.NewStore(store)

	// Create a valid user token
	token, _ := auth.CreateUserToken(testUser.Name, time.Hour)

	// Connect to amplifier
	log.Println("Connecting to amplifier")
	authenticatedConn, err := grpc.Dial(amp.AmplifierDefaultEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
		grpc.WithPerRPCCredentials(&auth.LoginCredentials{Token: token}),
	)
	if err != nil {
		log.Panicf("Unable to connect to amplifier on: %s\n%v", amp.AmplifierDefaultEndpoint, err)
	}
	anonymousConn, err := grpc.Dial(amp.AmplifierDefaultEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	if err != nil {
		log.Panicf("Unable to connect to amplifier on: %s\n%v", amp.AmplifierDefaultEndpoint, err)
	}
	log.Println("Connected to amplifier")

	// Init mail
	initMailServer()

	// Clients init
	functionClient = function.NewFunctionClient(authenticatedConn)
	statsClient = stats.NewStatsClient(authenticatedConn)
	stackClient = stack.NewStackServiceClient(authenticatedConn)
	topicClient = topic.NewTopicClient(authenticatedConn)
	serviceClient = service.NewServiceClient(authenticatedConn)
	logsClient = logs.NewLogsClient(authenticatedConn)
	accountClient = account.NewAccountClient(authenticatedConn)
	anonymousAccountClient = account.NewAccountClient(anonymousConn)

	// Start tests
	os.Exit(m.Run())
}
