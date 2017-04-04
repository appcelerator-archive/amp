package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	ctx           context.Context
	store         storage.Interface
	accountClient account.AccountClient
	accountStore  accounts.Interface
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Stores
	store = etcd.New([]string{etcd.DefaultEndpoint}, "amp", 5*time.Second)
	if err := store.Connect(); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", etcd.DefaultEndpoint, err)
	}
	accountStore = accounts.NewStore(store, configuration.RegistrationNone)

	// Connect to amplifier
	amplifierEndpoint := "amplifier" + configuration.DefaultPort
	log.Println("Connecting to amplifier")
	anonymousConn, err := grpc.Dial(amplifierEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
	)
	if err != nil {
		log.Panicf("Unable to connect to amplifier on: %s\n%v", amplifierEndpoint, err)
	}
	log.Println("Connected to amplifier")

	// Anonymous clients
	accountClient = account.NewAccountClient(anonymousConn)

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
	verificationToken, err := auth.CreateVerificationToken(user.Name)
	assert.NoError(t, err)

	// Verify
	_, err = accountClient.Verify(ctx, &account.VerificationRequest{Token: verificationToken})
	assert.NoError(t, err)

	// Login
	header := metadata.MD{}
	_, err = accountClient.Login(ctx, &account.LogInRequest{Name: user.Name, Password: user.Password}, grpc.Header(&header))
	assert.NoError(t, err)

	// Extract token from header
	tokens := header[auth.TokenKey]
	assert.NotEmpty(t, tokens)
	token := tokens[0]
	assert.NotEmpty(t, token)

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
	tokens := header[auth.TokenKey]
	assert.NotEmpty(t, tokens)
	token := tokens[0]
	assert.NotEmpty(t, token)

	return metadata.NewContext(ctx, metadata.Pairs(auth.TokenKey, token))
}

func changeOrganizationMemberRole(userCtx context.Context, t *testing.T, org *account.CreateOrganizationRequest, user *account.SignUpRequest, role accounts.OrganizationRole) {
	_, err := accountClient.ChangeOrganizationMemberRole(userCtx, &account.ChangeOrganizationMemberRoleRequest{OrganizationName: org.Name, UserName: user.Name, Role: role})
	assert.NoError(t, err)
}
