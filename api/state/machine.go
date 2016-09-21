package state

import (
	"context"
	"fmt"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/data/storage"
	"path"
)

const statesRootKey = "states"

// RuleSet describe allowed state transitions
type RuleSet [][]bool

// Machine StateMachine is the state machine
type Machine struct {
	ruleSet RuleSet
	store   storage.Interface
}

// NewMachine return a new state machine
func NewMachine(ruleSet RuleSet, store storage.Interface) Machine {
	return Machine{ruleSet: ruleSet, store: store}
}

// canTransition return whether or not you can transition between states
func (s *Machine) canTransition(from int32, to int32) bool {
	return s.ruleSet[from][to]
}

// GetState get state
func (s *Machine) GetState(id string) (int32, error) {
	state := &State{}
	if err := runtime.Store.Get(context.Background(), path.Join(statesRootKey, id), state, true); err != nil {
		return -1, err
	}
	return state.Value, nil
}

// TransitionTo transitionTo
func (s *Machine) TransitionTo(id string, to int32) error {
	current, err := s.GetState(id)
	if err != nil {
		return err
	}
	if !s.canTransition(current, to) {
		return fmt.Errorf("Cannot transition to state %v", to)
	}
	expect := &State{Value: current}
	update := &State{Value: to}
	if err = runtime.Store.CompareAndSet(context.Background(), path.Join(statesRootKey, id), expect, update); err != nil {
		return err
	}
	return nil
}

// Is is
func (s *Machine) Is(id string, expected int32) (bool, error) {
	state, err := s.GetState(id)
	if err != nil {
		return false, err
	}
	return state == expected, nil
}

// CreateState createstate
func (s *Machine) CreateState(id string, initial int32) error {
	state := &State{Value: initial}
	if err := runtime.Store.Create(context.Background(), path.Join(statesRootKey, id), state, nil, 0); err != nil {
		return err
	}
	return nil
}

// DeleteState deleteState
func (s *Machine) DeleteState(id string) error {
	return runtime.Store.Delete(context.Background(), path.Join(statesRootKey, id), false, nil)
}
