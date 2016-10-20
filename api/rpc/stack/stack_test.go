package stack_test

import (
	"os"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/api/server"
	"github.com/appcelerator/amp/api/state"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	example1 = `
services:
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

	example2 = `
services:
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
        publish_port: 3002
        internal_port: 3000`
)

var (
	client stack.StackServiceClient
	ctx    context.Context
)

func TestMain(m *testing.M) {
	_, conn := server.StartTestServer()
	ctx = context.Background()
	client = stack.NewStackServiceClient(conn)
	os.Exit(m.Run())
}

//Test two stacks life cycle in the same time
func TestShouldManageStackLifeCycleSuccessfully(t *testing.T) {
	//Start stack essai1
	//name := fmt.Sprintf("test-%d", time.Now().Unix())
	name := "stacktest"
	//Start stack test
	t.Log("start stack " + name)
	rUp, errUp := client.Up(ctx, &stack.UpRequest{StackName: name, Stackfile: example1})
	if errUp != nil {
		t.Fatal(errUp)
	}
	assert.NotEmpty(t, rUp.StackId, "Stack test1 StackId should not be empty")
	time.Sleep(3 * time.Second)
	//verifyusing ls
	t.Log("perform stack ls")
	listRequest := stack.ListRequest{}
	_, errls := client.List(ctx, &listRequest)
	if errls != nil {
		t.Fatal(errls)
	}
	//Prepare requests
	stackRequest := stack.StackRequest{
		StackIdent: rUp.StackId,
	}
	//Stop stack test
	time.Sleep(1 * time.Second)
	t.Log("stop stack " + name)
	rStop, errStop := client.Stop(ctx, &stackRequest)
	if errStop != nil {
		t.Fatal(errStop)
	}
	assert.NotEmpty(t, rStop.StackId, "Stack test StackId should not be empty")
	//Restart stack test
	time.Sleep(3 * time.Second)
	t.Log("restart stack " + name)
	rRestart, errRestart := client.Start(ctx, &stackRequest)
	if errRestart != nil {
		t.Fatal(errRestart)
	}
	assert.NotEmpty(t, rRestart.StackId, "Stack test StackId should not be empty")
	time.Sleep(3 * time.Second)
	t.Log("remove stack " + name)
	//Remove stack test
	removeRequest := stack.RemoveRequest{
		StackIdent: rUp.StackId,
		Force:      true,
	}
	rRemove, errRemove := client.Remove(ctx, &removeRequest)
	if errRemove != nil {
		t.Fatal(errRemove)
	}
	assert.NotEmpty(t, rRemove.StackId, "Stack test StackId should not be empty")
}

func TestTransitionsFromStopped(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)

	id := stringid.GenerateNonCryptoID()
	machine.CreateState(id, stack.StackState_Stopped.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Stopped.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Stopped.String())
	assert.NoError(t, machine.TransitionTo(id, stack.StackState_Starting.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Stopped.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Running.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Stopped.String())
	assert.NoError(t, machine.TransitionTo(id, stack.StackState_Redeploying.String()))
	machine.DeleteState(id)
}

func TestTransitionsFromStarting(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)
	id := stringid.GenerateNonCryptoID()

	machine.CreateState(id, stack.StackState_Starting.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Stopped.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Starting.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Starting.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Starting.String())
	assert.NoError(t, machine.TransitionTo(id, stack.StackState_Running.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Starting.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Redeploying.String()))
	machine.DeleteState(id)
}

func TestTransitionsFromRunning(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)
	id := stringid.GenerateNonCryptoID()

	machine.CreateState(id, stack.StackState_Running.String())
	assert.NoError(t, machine.TransitionTo(id, stack.StackState_Stopped.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Running.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Starting.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Running.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Running.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Running.String())
	assert.NoError(t, machine.TransitionTo(id, stack.StackState_Redeploying.String()))
	machine.DeleteState(id)
}

func TestTransitionsFromRedeploying(t *testing.T) {
	machine := state.NewMachine(stack.StackRuleSet, runtime.Store)
	id := stringid.GenerateNonCryptoID()

	machine.CreateState(id, stack.StackState_Redeploying.String())
	assert.NoError(t, machine.TransitionTo(id, stack.StackState_Stopped.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Redeploying.String())
	assert.NoError(t, machine.TransitionTo(id, stack.StackState_Starting.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Redeploying.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Running.String()))
	machine.DeleteState(id)

	machine.CreateState(id, stack.StackState_Redeploying.String())
	assert.Error(t, machine.TransitionTo(id, stack.StackState_Redeploying.String()))
	machine.DeleteState(id)
}
