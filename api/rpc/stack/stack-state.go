package stack

import (
	"fmt"
	"sync"
)

// State holds the current stack state, and has methods to get and
// set the state.
type State struct {
	sync.Mutex
	starting    bool
	running     bool
	stopped     bool
	redeploying bool
}

// NewState creates a default state object.
func NewState() *State {
	return &State{
		stopped: true,
	}
}

// IsStarting returns whether the starting flag is set.
func (s *State) IsStarting() bool {
	s.Lock()
	res := s.starting
	s.Unlock()
	return res
}

// IsRunning returns whether the running flag is set.
func (s *State) IsRunning() bool {
	s.Lock()
	res := s.running
	s.Unlock()
	return res
}

// IsStopped returns whether the stopped flag is set.
func (s *State) IsStopped() bool {
	s.Lock()
	res := s.stopped
	s.Unlock()
	return res
}

// IsRedeploying returns whether the redeploying flag is set.
func (s *State) IsRedeploying() bool {
	s.Lock()
	res := s.redeploying
	s.Unlock()
	return res
}

// SetStarting sets the state to "starting".
func (s *State) SetStarting() error {
	if !(s.stopped || s.redeploying) {
		return fmt.Errorf("Cannot transition to starting state")
	}
	s.Lock()
	s.stopped = false
	s.redeploying = false
	s.starting = true
	s.Unlock()
	return nil
}

// SetRunning sets the state to "running".
func (s *State) SetRunning() error {
	if !s.starting {
		return fmt.Errorf("Cannot transition to running state")
	}
	s.Lock()
	s.starting = false
	s.running = true
	s.Unlock()
	return nil
}

// SetStopped sets the state to "stopped".
func (s *State) SetStopped() error {
	if !(s.running || s.redeploying) {
		return fmt.Errorf("Cannot transition to stopped state")
	}
	s.Lock()
	s.running = false
	s.redeploying = false
	s.stopped = true
	s.Unlock()
	return nil
}

// SetRedeploying sets the state to "redeploying".
func (s *State) SetRedeploying() error {
	if !(s.running || s.stopped) {
		return fmt.Errorf("Cannot transition to redeploying state")
	}
	s.Lock()
	s.running = false
	s.stopped = false
	s.redeploying = true
	s.Unlock()
	return nil
}
