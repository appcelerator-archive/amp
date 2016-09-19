package stack

import (
	"fmt"
	"github.com/appcelerator/amp/api/runtime"
	"golang.org/x/net/context"
	"path"
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
func (s *Stack) IsStopped() (bool, error) { return s.testState(StackState_Stopped) }

// IsStarting returns true if the state is Starting
func (s *Stack) IsStarting() (bool, error) { return s.testState(StackState_Starting) }

// IsRunning returns true if the state is Running
func (s *Stack) IsRunning() (bool, error) { return s.testState(StackState_Running) }

// IsRedeploying returns true if the state is Redeploying
func (s *Stack) IsRedeploying() (bool, error) { return s.testState(StackState_Redeploying) }

// SetStopped sets the state to "stopped".
func (s *Stack) SetStopped() error { return s.transitionTo(StackState_Stopped) }

// SetStarting sets the state to "starting".
func (s *Stack) SetStarting() error { return s.transitionTo(StackState_Starting) }

// SetRunning sets the state to "running".
func (s *Stack) SetRunning() error { return s.transitionTo(StackState_Running) }

// SetRedeploying sets the state to "redeploying".
func (s *Stack) SetRedeploying() error { return s.transitionTo(StackState_Redeploying) }

func (s *Stack) getState() (StackState, error) {
	state := &State{}
	err := runtime.Store.Get(context.Background(), path.Join("stacks", s.Id, "state"), state, true)
	if err != nil {
		return -1, err
	}
	return state.Value, nil
}

func (s *Stack) transitionTo(to StackState) error {
	current, err := s.getState()
	if err != nil {
		return err
	}
	if !stackStateMachine.CanTransition(int32(current), int32(to)) {
		return fmt.Errorf("Cannot transition to state %s", to.String())
	}
	expect := &State{Value: current}
	update := &State{Value: to}
	err = runtime.Store.CompareAndSet(context.Background(), path.Join(stackRootKey, "/", s.Id, "/state"), expect, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Stack) testState(expected StackState) (bool, error) {
	state, err := s.getState()
	if err != nil {
		return false, err
	}
	return state == expected, nil
}
