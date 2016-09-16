package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransitionsFromStopped(t *testing.T) {
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Stopped), int32(StackState_Stopped)))
	assert.True(t, stackStateMachine.CanTransition(int32(StackState_Stopped), int32(StackState_Starting)))
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Stopped), int32(StackState_Running)))
	assert.True(t, stackStateMachine.CanTransition(int32(StackState_Stopped), int32(StackState_Redeploying)))
}

func TestTransitionsFromStarting(t *testing.T) {
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Starting), int32(StackState_Stopped)))
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Starting), int32(StackState_Starting)))
	assert.True(t, stackStateMachine.CanTransition(int32(StackState_Starting), int32(StackState_Running)))
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Starting), int32(StackState_Redeploying)))
}

func TestTransitionsFromRunning(t *testing.T) {
	assert.True(t, stackStateMachine.CanTransition(int32(StackState_Running), int32(StackState_Stopped)))
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Running), int32(StackState_Starting)))
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Running), int32(StackState_Running)))
	assert.True(t, stackStateMachine.CanTransition(int32(StackState_Running), int32(StackState_Redeploying)))
}

func TestTransitionsFromRedeploying(t *testing.T) {
	assert.True(t, stackStateMachine.CanTransition(int32(StackState_Redeploying), int32(StackState_Stopped)))
	assert.True(t, stackStateMachine.CanTransition(int32(StackState_Redeploying), int32(StackState_Starting)))
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Redeploying), int32(StackState_Running)))
	assert.False(t, stackStateMachine.CanTransition(int32(StackState_Redeploying), int32(StackState_Redeploying)))
}
