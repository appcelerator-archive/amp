package stack_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/api/server"
	"github.com/appcelerator/amp/api/state"
	"github.com/docker/docker/pkg/stringid"
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
	example                 = `
pinger:
  image: appcelerator/pinger
  replicas: 2
pingerExt1:
  image: appcelerator/pinger
  replicas: 2
  public:
    - name: www1
      protocol: tcp
      internal_port: 3000
pingerExt2:
  image: appcelerator/pinger
  replicas: 2
  public:
    - name: www2
      protocol: tcp
      publish_port: 3001
      internal_port: 3000`
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

func TestShouldUpStopRemoveStackSuccessfully(t *testing.T) {
	rUp, errUp := client.Up(ctx, &stack.UpRequest{Stackfile: example})
	if errUp != nil {
		t.Fatal(errUp)
	}
	assert.NotEmpty(t, rUp.StackId, "StackId should not be empty")
	fmt.Printf("Stack id = %s\n", rUp.StackId)
	stackRequest := stack.StackRequest{
		StackId: rUp.StackId,
	}
	rStop, errStop := client.Stop(ctx, &stackRequest)
	if errStop != nil {
		t.Fatal(errStop)
	}
	assert.NotEmpty(t, rStop.StackId, "StackId should not be empty")
	rRemove, errRemove := client.Remove(ctx, &stackRequest)
	if errRemove != nil {
		t.Fatal(errRemove)
	}
	assert.NotEmpty(t, rRemove.StackId, "StackId should not be empty")
}

func TestTransitionsFromStopped(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)

	id := stringid.GenerateNonCryptoID()
	machine.CreateState(id, int32(stack.StackState_Stopped))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Stopped)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Stopped))
	assert.NoError(t, machine.TransitionTo(id, int32(stack.StackState_Starting)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Stopped))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Running)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Stopped))
	assert.NoError(t, machine.TransitionTo(id, int32(stack.StackState_Redeploying)))
	machine.DeleteState(id)
}

func TestTransitionsFromStarting(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)
	id := stringid.GenerateNonCryptoID()

	machine.CreateState(id, int32(stack.StackState_Starting))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Stopped)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Starting))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Starting)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Starting))
	assert.NoError(t, machine.TransitionTo(id, int32(stack.StackState_Running)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Starting))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Redeploying)))
	machine.DeleteState(id)
}

func TestTransitionsFromRunning(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)
	id := stringid.GenerateNonCryptoID()

	machine.CreateState(id, int32(stack.StackState_Running))
	assert.NoError(t, machine.TransitionTo(id, int32(stack.StackState_Stopped)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Running))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Starting)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Running))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Running)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Running))
	assert.NoError(t, machine.TransitionTo(id, int32(stack.StackState_Redeploying)))
	machine.DeleteState(id)
}

func TestTransitionsFromRedeploying(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)
	id := stringid.GenerateNonCryptoID()

	machine.CreateState(id, int32(stack.StackState_Redeploying))
	assert.NoError(t, machine.TransitionTo(id, int32(stack.StackState_Stopped)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Redeploying))
	assert.NoError(t, machine.TransitionTo(id, int32(stack.StackState_Starting)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Redeploying))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Running)))
	machine.DeleteState(id)

	machine.CreateState(id, int32(stack.StackState_Redeploying))
	assert.Error(t, machine.TransitionTo(id, int32(stack.StackState_Redeploying)))
	machine.DeleteState(id)
}
