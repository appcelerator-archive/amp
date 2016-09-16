package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransitionsFromStopped(t *testing.T) {
	assert.False(t, stackStateMachine.CanTransition(Stopped, Stopped))
	assert.True(t, stackStateMachine.CanTransition(Stopped, Starting))
	assert.False(t, stackStateMachine.CanTransition(Stopped, Running))
	assert.True(t, stackStateMachine.CanTransition(Stopped, Redeploying))
}

func TestTransitionsFromStarting(t *testing.T) {
	assert.False(t, stackStateMachine.CanTransition(Starting, Stopped))
	assert.False(t, stackStateMachine.CanTransition(Starting, Starting))
	assert.True(t, stackStateMachine.CanTransition(Starting, Running))
	assert.False(t, stackStateMachine.CanTransition(Starting, Redeploying))
}

func TestTransitionsFromRunning(t *testing.T) {
	assert.True(t, stackStateMachine.CanTransition(Running, Stopped))
	assert.False(t, stackStateMachine.CanTransition(Running, Starting))
	assert.False(t, stackStateMachine.CanTransition(Running, Running))
	assert.True(t, stackStateMachine.CanTransition(Running, Redeploying))
}

func TestTransitionsFromRedeploying(t *testing.T) {
	assert.True(t, stackStateMachine.CanTransition(Redeploying, Stopped))
	assert.True(t, stackStateMachine.CanTransition(Redeploying, Starting))
	assert.False(t, stackStateMachine.CanTransition(Redeploying, Running))
	assert.False(t, stackStateMachine.CanTransition(Redeploying, Redeploying))
}
