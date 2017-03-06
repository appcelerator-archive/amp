package tests

import (
	"github.com/appcelerator/amp/api/authn"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/functions"
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
	accountStore   accounts.Interface
	functionStore  functions.Interface
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Stores
	store = etcd.New([]string{amp.EtcdDefaultEndpoint}, "amp")
	if err := store.Connect(5 * time.Second); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", amp.EtcdDefaultEndpoint, err)
	}
	accountStore = accounts.NewStore(store)
	functionStore = functions.NewStore(store)

	// Create a valid user token
	token, _ := authn.CreateLoginToken("default", "", time.Hour)

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
	verificationToken, err := authn.CreateVerificationToken(user.Name, time.Hour)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: verificationToken})
	assert.NoError(t, err)

	// Login
	header := metadata.MD{}
	_, err = accountClient.Login(ctx, &account.LogInRequest{Name: user.Name, Password: user.Password}, grpc.Header(&header))
	assert.NoError(t, err)

	// Extract token from header
	tokens := header[authn.TokenKey]
	assert.NotEmpty(t, tokens)
	token := tokens[0]
	assert.NotEmpty(t, token)

	return metadata.NewContext(ctx, metadata.Pairs(authn.TokenKey, token))
}

func createOrganization(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest) context.Context {
	// Create a user
	ownerCtx := createUser(t, owner)

	// CreateOrganization
	_, err := accountClient.CreateOrganization(ownerCtx, org)
	assert.NoError(t, err)

	return ownerCtx
}

func createAndAddUserToOrganization(ownerCtx context.Context, t *testing.T, org *account.CreateOrganizationRequest, user *account.SignUpRequest) context.Context {
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

func switchAccount(userCtx context.Context, t *testing.T, accountName string) context.Context {
	header := metadata.MD{}
	_, err := accountClient.Switch(userCtx, &account.SwitchRequest{Account: accountName}, grpc.Header(&header))
	assert.NoError(t, err)

	// Extract token from header
	tokens := header[authn.TokenKey]
	assert.NotEmpty(t, tokens)
	token := tokens[0]
	assert.NotEmpty(t, token)

	return metadata.NewContext(ctx, metadata.Pairs(authn.TokenKey, token))
}

func changeOrganizationMemberRole(userCtx context.Context, t *testing.T, org *account.CreateOrganizationRequest, user *account.SignUpRequest, role accounts.OrganizationRole) {
	_, err := accountClient.ChangeOrganizationMemberRole(userCtx, &account.ChangeOrganizationMemberRoleRequest{OrganizationName: org.Name, UserName: user.Name, Role: role})
	assert.NoError(t, err)
}
