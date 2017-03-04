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
	"github.com/appcelerator/amp/cmd/amp/cli"
	as "github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	// Stores
	store = etcd.New([]string{amp.EtcdDefaultEndpoint}, "amp")
	if err := store.Connect(5 * time.Second); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", amp.EtcdDefaultEndpoint, err)
	}
	accountStore = as.NewStore(store)

	// Create a valid user token
	token, _ := auth.CreateLoginToken("default", "", time.Hour)

	// Connect to amplifier
	log.Println("Connecting to amplifier")
	authenticatedConn, err := grpc.Dial(amp.AmplifierDefaultEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
		grpc.WithPerRPCCredentials(&cli.LoginCredentials{Token: token}),
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

	// Authenticated clients
	statsClient = stats.NewStatsClient(authenticatedConn)
	stackClient = stack.NewStackServiceClient(authenticatedConn)
	topicClient = topic.NewTopicClient(authenticatedConn)
	serviceClient = service.NewServiceClient(authenticatedConn)
	logsClient = logs.NewLogsClient(authenticatedConn)

	// Anonymous clients
	accountClient = account.NewAccountClient(anonymousConn)
	functionClient = function.NewFunctionClient(anonymousConn)

	// Start tests
	code := m.Run()

	// Tear down
	accountStore.Reset(ctx)

	os.Exit(code)
}

func createUser(t *testing.T, user *account.SignUpRequest) context.Context {
	// SignUp
	_, err := accountClient.SignUp(ctx, user)
	assert.NoError(t, err)

	// Create a verify token
	token, err := auth.CreateVerificationToken(user.Name, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: token})
	assert.NoError(t, err)

	// Create a login token
	token, err = auth.CreateLoginToken(user.Name, "", time.Hour)
	return metadata.NewContext(ctx, metadata.Pairs(auth.TokenKey, token))
}

func createOrganization(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest) context.Context {
	// Create a user
	ownerCtx := createUser(t, owner)

	// CreateOrganization
	_, err := accountClient.CreateOrganization(ownerCtx, org)
	assert.NoError(t, err)

	return ownerCtx
}

func addUserToOrganization(t *testing.T, org *account.CreateOrganizationRequest, ownerCtx context.Context, user *account.SignUpRequest) context.Context {
	// Create a user
	userCtx := createUser(t, user)

	// AddUserToOrganization
	_, err := accountClient.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: org.Name,
		UserName:         user.Name,
	})
	assert.NoError(t, err)
	return userCtx
}

func createTeam(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest, team *account.CreateTeamRequest) context.Context {
	// Create a user
	ownerCtx := createOrganization(t, org, owner)

	// CreateTeam
	_, err := accountClient.CreateTeam(ownerCtx, team)
	assert.NoError(t, err)

	return ownerCtx
}
