package stack_test

import (
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/api/server"
	"github.com/appcelerator/amp/api/state"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
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
