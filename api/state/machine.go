package state

import (
	"context"
	"fmt"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/data/storage"
	"path"
)

// RuleSet describe allowed state transitions
type RuleSet [][]bool

// StateMachine is the state machine
type Machine struct {
	ruleSet RuleSet
	store   storage.Interface
}

// NewStateMachine return a new state machine
func NewMachine(ruleSet RuleSet, store storage.Interface) Machine {
	return Machine{ruleSet: ruleSet, store: store}
}

// canTransition return whether or not you can transition between states
func (s *Machine) canTransition(from int32, to int32) bool {
	return s.ruleSet[from][to]
}

func (s *Machine) getState(id string) (int32, error) {
	state := &State{}
	if err := runtime.Store.Get(context.Background(), path.Join("states", id), state, true); err != nil {
		return -1, err
	}
	return state.Value, nil
}

func (s *Machine) TransitionTo(id string, to int32) error {
	current, err := s.getState(id)
	if err != nil {
		return err
	}
	if !s.canTransition(current, to) {
		return fmt.Errorf("Cannot transition to state %s", to)
	}
	expect := &State{Value: current}
	update := &State{Value: to}
	if err = runtime.Store.CompareAndSet(context.Background(), path.Join("states", id), expect, update); err != nil {
		return err
	}
	return nil
}

func (s *Machine) Is(id string, expected int32) (bool, error) {
	state, err := s.getState(id)
	if err != nil {
		return false, err
	}
	return state == expected, nil
}

func (s *Machine) CreateState(id string, initial int32) error {
	state := &State{Value: initial}
	if err := runtime.Store.Create(context.Background(), path.Join("states", id), state, nil, 0); err != nil {
		return err
	}
	return nil
}

func (s *Machine) DeleteState(id string) error {
	return runtime.Store.Delete(context.Background(), path.Join("states", id), false, nil)
}
