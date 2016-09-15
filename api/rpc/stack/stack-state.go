package stack

import (
	"fmt"
)

type State int

const (
	Stopped State = iota
	Starting
	Running
	Redeploying
)

type stateMachine [][]bool

var stackStateMachine = stateMachine{
	//   | Stopped   | Starting  | Running   | Redeploying
	[]bool{false /**/, true /* */, false /**/, true /* */}, // Stopped (initial state)
	[]bool{false /**/, false /**/, true /* */, false /**/}, // Starting
	[]bool{true /* */, false /**/, false /**/, true /* */}, // Running
	[]bool{true /* */, true /* */, false /**/, false /**/}, // Redeploying
}

// NewState creates a default state object.
func NewState() State {
	return Stopped
}

// IsStarting returns whether the starting flag is set.
func (s *State) IsStarting() bool {
	return *s == Starting
}

// IsRunning returns whether the running flag is set.
func (s *State) IsRunning() bool {
	return *s == Running
}

// IsStopped returns whether the stopped flag is set.
func (s *State) IsStopped() bool {
	return *s == Stopped
}

// IsRedeploying returns whether the redeploying flag is set.
func (s *State) IsRedeploying() bool {
	return *s == Redeploying
}

// SetStarting sets the state to "starting".
func (s *State) SetStarting() error {
	if !stackStateMachine[*s][Starting] {
		return fmt.Errorf("Cannot transition to Starting state")
	}
	*s = Starting
	return nil
}

// SetRunning sets the state to "running".
func (s *State) SetRunning() error {
	if !stackStateMachine[*s][Running] {
		return fmt.Errorf("Cannot transition to Running state")
	}
	*s = Running
	return nil
}

// SetStopped sets the state to "stopped".
func (s *State) SetStopped() error {
	if !stackStateMachine[*s][Stopped] {
		return fmt.Errorf("Cannot transition to Stopped state")
	}
	*s = Stopped
	return nil
}

// SetRedeploying sets the state to "redeploying".
func (s *State) SetRedeploying() error {
	if !stackStateMachine[*s][Redeploying] {
		return fmt.Errorf("Cannot transition to Redeploying state")
	}
	*s = Redeploying
	return nil
}
