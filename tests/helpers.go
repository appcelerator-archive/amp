package helpers

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/cluster/constants"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/docker/docker/pkg/stringid"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
)

const testAmplifierConfig = "test.amplifier.yml"

// AmplifierConnection returns a grpc connection to amplifier
func AmplifierConnection() (*grpc.ClientConn, error) {
	amplifierEndpoint := "amplifier" + configuration.DefaultPort
	log.Println("Connecting to amplifier")
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.Dial(amplifierEndpoint,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to amplifier on: %store\n%v", amplifierEndpoint, err)
	}
	log.Println("Connected to amplifier")
	return conn, nil
}

// Helper is a test helper
type Helper struct {
	accounts   account.AccountClient
	logs       logs.LogsClient
	resources  resource.ResourceClient
	stacks     stack.StackClient
	tokens     *auth.Tokens
	suPassword string
}

// New returns a new test helper
func New() (*Helper, error) {
	conn, err := AmplifierConnection()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile("../../" + testAmplifierConfig)
	if err != nil {
		return nil, err
	}

	cfg := configuration.Configuration{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	h := &Helper{
		accounts:   account.NewAccountClient(conn),
		logs:       logs.NewLogsClient(conn),
		stacks:     stack.NewStackClient(conn),
		resources:  resource.NewResourceClient(conn),
		tokens:     auth.New(cfg.JWTSecretKey),
		suPassword: cfg.SUPassword,
	}
	return h, nil
}

// Accounts returns an account client
func (h *Helper) Accounts() account.AccountClient {
	return h.accounts
}

// Logs returns a log client
func (h *Helper) Logs() logs.LogsClient {
	return h.logs
}

// Stacks returns a stack client
func (h *Helper) Stacks() stack.StackClient {
	return h.stacks
}

// Resources returns a resource client
func (h *Helper) Resources() resource.ResourceClient {
	return h.resources
}

// Tokens returns a token manager instance
func (h *Helper) Tokens() *auth.Tokens {
	return h.tokens
}

// Login sign-up, verify and log-in a random user and returns the associated credentials
func (h *Helper) Login() (metadata.MD, error) {
	randomUser := h.RandomUser()
	conn, err := AmplifierConnection()
	if err != nil {
		return nil, err
	}
	client := account.NewAccountClient(conn)
	ctx := context.Background()

	// SignUp
	_, err = client.SignUp(ctx, &randomUser)
	if err != nil {
		return nil, fmt.Errorf("SignUp error: %v", err)
	}

	// Login
	header := metadata.MD{}
	_, err = client.Login(ctx, &account.LogInRequest{Name: randomUser.Name, Password: randomUser.Password}, grpc.Header(&header))
	if err != nil {
		return nil, fmt.Errorf("Login error: %v", err)
	}

	// Extract token from header
	tokens := header[auth.TokenKey]
	if len(tokens) == 0 {
		return nil, errors.New("No token in response header")
	}
	token := tokens[0]
	return metadata.Pairs(auth.AuthorizationHeader, auth.ForgeAuthorizationHeader(token)), nil
}

// SuperLogin logs in as super user and returns the associated credentials
func (h *Helper) SuperLogin() (context.Context, error) {
	su := h.SuperUser()
	conn, err := AmplifierConnection()
	if err != nil {
		return nil, err
	}
	client := account.NewAccountClient(conn)

	// Login
	header := metadata.MD{}
	_, err = client.Login(context.Background(), &account.LogInRequest{Name: su.Name, Password: su.Password}, grpc.Header(&header))
	if err != nil {
		return nil, fmt.Errorf("Login error: %v", err)
	}

	// Extract token from header
	tokens := header[auth.TokenKey]
	if len(tokens) == 0 {
		return nil, errors.New("No token in response header")
	}
	token := tokens[0]
	credentials := metadata.Pairs(auth.AuthorizationHeader, auth.ForgeAuthorizationHeader(token))
	return metadata.NewOutgoingContext(context.Background(), credentials), nil
}

// SuperUser returns the super user SignUpRequest
func (h *Helper) SuperUser() account.SignUpRequest {
	return account.SignUpRequest{
		Name:     "su",
		Password: h.suPassword,
		Email:    "super@user.amp",
	}
}

// RandomUser returns a random user SignUpRequest
func (h *Helper) RandomUser() account.SignUpRequest {
	id := stringid.GenerateNonCryptoID()
	return account.SignUpRequest{
		Name:     id,
		Password: "password",
		Email:    id + "@user.email",
	}
}

// RandomOrg returns a random organization CreateOrganizationRequest
func (h *Helper) RandomOrg() account.CreateOrganizationRequest {
	id := stringid.GenerateNonCryptoID()
	return account.CreateOrganizationRequest{
		Name:  id,
		Email: id + "@org.email",
	}
}

// RandomTeam returns a random team CreateTeamRequest
func (h *Helper) RandomTeam(org string) account.CreateTeamRequest {
	id := stringid.GenerateNonCryptoID()
	return account.CreateTeamRequest{
		OrganizationName: org,
		TeamName:         id,
	}
}

// DefaultOrg returns the default organization CreateOrganizationRequest
func (h *Helper) DefaultOrg() account.CreateOrganizationRequest {
	return account.CreateOrganizationRequest{
		Name:  accounts.DefaultOrganization,
		Email: accounts.DefaultOrganizationEmail,
	}
}

// DeployStack deploys a stack with given name and file location
func (h *Helper) DeployStack(ctx context.Context, name string, composeFile string) (string, error) {
	contents, err := ioutil.ReadFile(composeFile)
	if err != nil {
		return "", err
	}
	request := &stack.DeployRequest{
		Name:    name,
		Compose: contents,
	}
	reply, err := h.stacks.StackDeploy(ctx, request)
	if err != nil {
		return "", err
	}
	return reply.Id, err
}

// CreateUser signs up and logs in the given user
func (h *Helper) CreateUser(t *testing.T, user *account.SignUpRequest) context.Context {
	ctx := context.Background()

	// SignUp
	_, err := h.accounts.SignUp(ctx, user)
	assert.NoError(t, err)

	// Login
	header := metadata.MD{}
	_, err = h.accounts.Login(ctx, &account.LogInRequest{Name: user.Name, Password: user.Password}, grpc.Header(&header))
	assert.NoError(t, err)

	// Extract token from header
	tokens := header[auth.TokenKey]
	assert.NotEmpty(t, tokens)
	token := tokens[0]
	assert.NotEmpty(t, token)

	return metadata.NewOutgoingContext(ctx, metadata.Pairs(auth.AuthorizationHeader, auth.ForgeAuthorizationHeader(token)))
}

// CreateOrganization signs up and logs in the given user and creates an organization
func (h *Helper) CreateOrganization(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest) context.Context {
	// Create a user
	ownerCtx := h.CreateUser(t, owner)

	// CreateOrganization
	_, err := h.accounts.CreateOrganization(ownerCtx, org)
	assert.NoError(t, err)

	return ownerCtx
}

// CreateAndAddUserToOrganization signs up and logs in the given user and adds them to the given organization
func (h *Helper) CreateAndAddUserToOrganization(ownerCtx context.Context, t *testing.T, org *account.CreateOrganizationRequest, user *account.SignUpRequest) context.Context {
	// Create a user
	userCtx := h.CreateUser(t, user)

	// AddUserToOrganization
	_, err := h.accounts.AddUserToOrganization(ownerCtx, &account.AddUserToOrganizationRequest{
		OrganizationName: org.Name,
		UserName:         user.Name,
	})
	assert.NoError(t, err)
	return userCtx
}

// CreateTeam creates the given owner, creates the given organization and creates a team in this organization
func (h *Helper) CreateTeam(t *testing.T, org *account.CreateOrganizationRequest, owner *account.SignUpRequest, team *account.CreateTeamRequest) context.Context {
	// Create a user
	ownerCtx := h.CreateUser(t, owner)

	// CreateTeam
	_, err := h.accounts.CreateTeam(ownerCtx, team)
	assert.NoError(t, err)

	return ownerCtx
}

func (h *Helper) AddUserToTeam(ownerCtx context.Context, t *testing.T, team *account.CreateTeamRequest, user *account.SignUpRequest) {
	_, err := h.Accounts().AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: team.OrganizationName,
		TeamName:         team.TeamName,
		UserName:         user.Name,
	})
	assert.NoError(t, err)
}

func (h *Helper) RemoveUserFromTeam(ownerCtx context.Context, t *testing.T, team *account.CreateTeamRequest, user *account.SignUpRequest) {
	_, err := h.Accounts().RemoveUserFromTeam(ownerCtx, &account.RemoveUserFromTeamRequest{
		OrganizationName: team.OrganizationName,
		TeamName:         team.TeamName,
		UserName:         user.Name,
	})
	assert.NoError(t, err)
}

// Switch switches from the given user context to the given account name
func (h *Helper) Switch(userCtx context.Context, t *testing.T, accountName string) context.Context {
	header := metadata.MD{}
	_, err := h.accounts.Switch(userCtx, &account.SwitchRequest{Account: accountName}, grpc.Header(&header))
	assert.NoError(t, err)

	// Extract token from header
	tokens := header[auth.TokenKey]
	assert.NotEmpty(t, tokens)
	token := tokens[0]
	assert.NotEmpty(t, token)

	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(auth.AuthorizationHeader, auth.ForgeAuthorizationHeader(token)))
}

// Log Producer constants
var (
	TestMessage            = "test message "
	TestContainerID        = stringid.GenerateNonCryptoID()
	TestContainerName      = "testcontainer"
	TestContainerShortName = "testcontainershortname"
	TestContainerState     = "testcontainerstate"
	TestServiceID          = stringid.GenerateNonCryptoID()
	TestServiceName        = "testservice"
	TestStackID            = stringid.GenerateNonCryptoID()
	TestStackName          = "teststack"
	TestNodeID             = stringid.GenerateNonCryptoID()
	TestTaskID             = stringid.GenerateNonCryptoID()
)

// LogProducer is a test log producer
type LogProducer struct {
	ns              *ns.NatsStreaming
	asyncProduction int32
	counter         int64
	h               *Helper
}

// NewLogProducer instantiates a new log producer
func NewLogProducer(h *Helper) *LogProducer {
	lp := &LogProducer{
		ns:      ns.NewClient(ns.DefaultURL, ns.ClusterID, stringid.GenerateNonCryptoID(), 60*time.Second),
		counter: 0,
		h:       h,
	}
	if err := lp.ns.Connect(); err != nil {
		log.Fatalln("Cannot connect to NATS", err)
	}
	go func(lp *LogProducer) {
		for {
			time.Sleep(50 * time.Millisecond)
			if lp.asyncProduction > 0 {
				if err := lp.produce(logs.NumberOfEntries); err != nil {
					log.Println("error producing async messages", err)
				}
			}
		}
	}(lp)
	return lp
}

func (lp *LogProducer) buildLogEntry(infrastructure bool) *logs.LogEntry {
	atomic.AddInt64(&lp.counter, 1)
	entry := &logs.LogEntry{
		Timestamp:          time.Now().UTC().Format(time.RFC3339Nano),
		ContainerId:        TestContainerID,
		ContainerName:      TestContainerName,
		ContainerShortName: TestContainerShortName,
		ContainerState:     TestContainerState,
		ServiceName:        TestServiceName,
		ServiceId:          TestServiceID,
		TaskId:             TestTaskID,
		StackName:          TestStackName,
		StackId:            TestStackID,
		NodeId:             TestNodeID,
		TimeId:             fmt.Sprintf("%016X", lp.counter),
		Labels:             make(map[string]string),
		Msg:                TestMessage + fmt.Sprintf("%016X", lp.counter),
	}
	if infrastructure {
		entry.Labels[constants.LabelKeyRole] = "infra"
	}
	return entry
}

func (lp *LogProducer) produce(howMany int) error {
	entries := logs.GetReply{}
	for i := 0; i < howMany; i++ {
		// User log entry
		user := lp.buildLogEntry(false)
		entries.Entries = append(entries.Entries, user)

		// Infrastructure log entry
		infra := lp.buildLogEntry(true)
		entries.Entries = append(entries.Entries, infra)
	}
	message, err := proto.Marshal(&entries)
	if err != nil {
		return err
	}
	if err := lp.ns.GetClient().Publish(ns.LogsSubject, message); err != nil {
		return err
	}
	return nil
}

// StartAsyncProducer starts producing log messages in the background
func (lp *LogProducer) StartAsyncProducer() {
	atomic.CompareAndSwapInt32(&lp.asyncProduction, 0, 1)
}

// StopAsyncProducer stops producing log messages in the background
func (lp *LogProducer) StopAsyncProducer() {
	atomic.CompareAndSwapInt32(&lp.asyncProduction, 1, 0)
}

// PopulateLogs populate elasticsearch with log entries, by producing test messages and making sure their got indexed
func (lp *LogProducer) PopulateLogs() error {
	// Connect to amplifier
	conn, err := AmplifierConnection()
	if err != nil {
		return err
	}
	client := logs.NewLogsClient(conn)

	// Login context
	credentials, err := lp.h.Login()
	if err != nil {
		return err
	}
	ctx := metadata.NewOutgoingContext(context.Background(), credentials)

	// Populate logs
	if err := lp.produce(logs.NumberOfEntries); err != nil {
		return err
	}

	// Wait for them to be indexed in elasticsearch
	for {
		time.Sleep(1 * time.Second)
		r, err := client.LogsGet(ctx, &logs.GetRequest{Service: TestServiceID})
		if err != nil {
			log.Println(err)
			continue
		}
		if len(r.Entries) == logs.NumberOfEntries {
			break
		}
	}
	return nil
}
