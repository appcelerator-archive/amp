package stack_test

import (
	"fmt"
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
	name1 := fmt.Sprintf("test1-%d", time.Now().Unix())
	name2 := fmt.Sprintf("test2-%d", time.Now().Unix())
	//Start stack test1
	t.Log("start stack " + name1)
	rUp1, errUp1 := client.Up(ctx, &stack.UpRequest{StackName: name1, Stackfile: example1})
	if errUp1 != nil {
		t.Fatal(errUp1)
	}
	//Start stack test2
	t.Log("start stack " + name2)
	rUp2, errUp2 := client.Up(ctx, &stack.UpRequest{StackName: name2, Stackfile: example2})
	if errUp2 != nil {
		t.Fatal(errUp2)
	}
	assert.NotEmpty(t, rUp1.StackId, "Stack test1 StackId should not be empty")
	assert.NotEmpty(t, rUp2.StackId, "Stack test2 StackId should not be empty")
	time.Sleep(3 * time.Second)
	//verifyusing ls
	t.Log("perform stack ls")
	listRequest := stack.ListRequest{}
	_, errls := client.List(ctx, &listRequest)
	if errls != nil {
		t.Fatal(errls)
	}
	//Prepare requests
	stackRequest1 := stack.StackRequest{
		StackIdent: rUp1.StackId,
	}
	stackRequest2 := stack.StackRequest{
		StackIdent: rUp2.StackId,
	}
	//Stop stack test1
	t.Log("stop stack " + name1)
	rStop1, errStop1 := client.Stop(ctx, &stackRequest1)
	if errStop1 != nil {
		t.Fatal(errStop1)
	}
	assert.NotEmpty(t, rStop1.StackId, "Stack test1 StackId should not be empty")
	//Restart stack test1
	time.Sleep(1 * time.Second)
	t.Log("restart stack " + name1)
	rRestart1, errRestart1 := client.Start(ctx, &stackRequest1)
	if errRestart1 != nil {
		t.Fatal(errRestart1)
	}
	assert.NotEmpty(t, rRestart1.StackId, "Stack test1 StackId should not be empty")
	time.Sleep(1 * time.Second)
	//Stop again stack test1
	t.Log("stop again stack " + name1)
	rStop12, errStop12 := client.Stop(ctx, &stackRequest1)
	if errStop12 != nil {
		t.Fatal(errStop12)
	}
	assert.NotEmpty(t, rStop12.StackId, "Stack test1 StackId should not be empty")
	t.Log("remove stack " + name1)
	//Remove stack test1
	removeRequest1 := stack.RemoveRequest{
		StackIdent: rUp1.StackId,
		Force:      false,
	}
	rRemove1, errRemove1 := client.Remove(ctx, &removeRequest1)
	if errRemove1 != nil {
		t.Fatal(errRemove1)
	}
	assert.NotEmpty(t, rRemove1.StackId, "Stack test1 StackId should not be empty")
	//Stop stack test2
	t.Log("stop stack " + name2)
	rStop2, errStop2 := client.Stop(ctx, &stackRequest2)
	if errStop2 != nil {
		t.Fatal(errStop2)
	}
	assert.NotEmpty(t, rStop2.StackId, "Stack test2 StackId should not be empty")
	//Remove stack test2
	t.Log("remove stack " + name2)
	removeRequest2 := stack.RemoveRequest{
		StackIdent: rUp2.StackId,
		Force:      false,
	}
	rRemove2, errRemove2 := client.Remove(ctx, &removeRequest2)
	if errRemove2 != nil {
		t.Fatal(errRemove2)
	}
	assert.NotEmpty(t, rRemove2.StackId, "Stack test2 StackId should not be empty")
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
