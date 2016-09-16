package stack

import (
	"fmt"
	"github.com/appcelerator/amp/api/runtime"
	"golang.org/x/net/context"
	"path"
)

const (
	// Stopped (initial state)
	Stopped = iota
	// Starting state
	Starting
	// Running  state
	Running
	// Redeploying  state
	Redeploying
)

var stackRuleSet = RuleSet{
	//   | Stopped   | Starting  | Running   | Redeploying
	[]bool{false /**/, true /* */, false /**/, true /* */}, // Stopped (initial state)
	[]bool{false /**/, false /**/, true /* */, false /**/}, // Starting
	[]bool{true /* */, false /**/, false /**/, true /* */}, // Running
	[]bool{true /* */, true /* */, false /**/, false /**/}, // Redeploying
}

var stackStateMachine = NewStateMachine(stackRuleSet)

// IsStopped returns true if the state is Stopped
func (s *Stack) IsStopped() (bool, error) { return s.testState(Stopped) }

// IsStarting returns true if the state is Starting
func (s *Stack) IsStarting() (bool, error) { return s.testState(Starting) }

// IsRunning returns true if the state is Running
func (s *Stack) IsRunning() (bool, error) { return s.testState(Running) }

// IsRedeploying returns true if the state is Redeploying
func (s *Stack) IsRedeploying() (bool, error) { return s.testState(Redeploying) }

// SetStopped sets the state to "stopped".
func (s *Stack) SetStopped() error { return s.transitionTo(Stopped) }

// SetStarting sets the state to "starting".
func (s *Stack) SetStarting() error { return s.transitionTo(Starting) }

// SetRunning sets the state to "running".
func (s *Stack) SetRunning() error { return s.transitionTo(Running) }

// SetRedeploying sets the state to "redeploying".
func (s *Stack) SetRedeploying() error { return s.transitionTo(Redeploying) }

func (s *Stack) getState() (int32, error) {
	stack := &Stack{}
	err := runtime.Store.Get(context.Background(), path.Join("stacks", s.Id), stack, false)
	if err != nil {
		return -1, err
	}
	return stack.State, nil
}

func (s *Stack) setState(state int32) error {
	stack := &Stack{}
	err := runtime.Store.Get(context.Background(), path.Join("stacks", s.Id), stack, false)
	if err != nil {
		return err
	}
	stack.State = state
	err = runtime.Store.Update(context.Background(), path.Join("stacks", s.Id), stack, int64(0))
	if err != nil {
		return err
	}
	return nil
}

func (s *Stack) transitionTo(to int32) error {
	current, err := s.getState()
	if err != nil {
		return err
	}
	if !stackStateMachine.CanTransition(current, to) {
		return fmt.Errorf("Cannot transition to %d state", to)
	}
	s.setState(to)
	if err != nil {
		return err
	}
	return nil
}

func (s *Stack) testState(expected int32) (bool, error) {
	state, err := s.getState()
	if err != nil {
		return false, err
	}
	return state == expected, nil
}
