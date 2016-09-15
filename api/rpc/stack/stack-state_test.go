package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldCreateAValidInitilaState(t *testing.T) {
	state := NewState()
	assert.True(t, state.IsStopped(), "State should be stopped initially")
	assert.False(t, state.IsStarting(), "State should be stopped initially")
	assert.False(t, state.IsRunning(), "State should be stopped initially")
	assert.False(t, state.IsRedeploying(), "State should be stopped initially")
}

func TestTransitionsFromStopped(t *testing.T) {
	// KO
	assert.Error(t, NewState().SetStopped(), "transitioning from stopped to stopped should not be possible")
	assert.Error(t, NewState().SetRunning(), "transitioning from stopped to running should not be possible")

	// OK
	assert.NoError(t, NewState().SetRedeploying(), "transitioning from stopped to redeploying should be possible")
	assert.NoError(t, NewState().SetStarting(), "transitioning from stopped to starting should be possible")
}

func TestTransitionsFromStarting(t *testing.T) {
	// KO
	state := NewState()
	state.SetStarting()
	assert.Error(t, state.SetStopped(), "transitioning from starting to stopped should not be possible")
	state = NewState()
	state.SetStarting()
	assert.Error(t, state.SetStarting(), "transitioning from starting to starting should not be possible")
	state = NewState()
	state.SetStarting()
	assert.Error(t, state.SetRedeploying(), "transitioning from starting to redeploying should not be possible")

	// OK
	state = NewState()
	state.SetStarting()
	assert.NoError(t, state.SetRunning(), "transitioning from starting to running should be possible")
}

func TestTransitionsFromRunning(t *testing.T) {
	// KO
	state := NewState()
	state.SetStarting()
	state.SetRunning()
	assert.Error(t, state.SetRunning(), "transitioning from running to running should not be possible")
	state = NewState()
	state.SetStarting()
	state.SetRunning()
	assert.Error(t, state.SetStarting(), "transitioning from running to starting should not be possible")

	// OK
	state = NewState()
	state.SetStarting()
	state.SetRunning()
	assert.NoError(t, state.SetRedeploying(), "transitioning from running to redeploying should be possible")
	state = NewState()
	state.SetStarting()
	state.SetRunning()
	assert.NoError(t, state.SetStopped(), "transitioning from running to stopped should be possible")
}

func TestTransitionsFromRedeploying(t *testing.T) {
	// KO
	state := NewState()
	state.SetStarting()
	state.SetRunning()
	state.SetRedeploying()
	assert.Error(t, state.SetRedeploying(), "transitioning from redeploying to redeploying should not be possible")
	state = NewState()
	state.SetStarting()
	state.SetRunning()
	state.SetRedeploying()
	assert.Error(t, state.SetRunning(), "transitioning from redeploying to running should not be possible")

	// OK
	state = NewState()
	state.SetStarting()
	state.SetRunning()
	state.SetRedeploying()
	assert.NoError(t, state.SetStopped(), "transitioning from redeploying to stopped should be possible")
	state = NewState()
	state.SetStarting()
	state.SetRunning()
	state.SetRedeploying()
	assert.NoError(t, state.SetStarting(), "transitioning from redeploying to starting should be possible")
}
