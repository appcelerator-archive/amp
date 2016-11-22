package state

import (
	"context"
	"fmt"
	"path"

	"github.com/appcelerator/amp/data/storage"
)

const statesRootKey = "states"

// RuleSet describe allowed state transitions
type RuleSet map[string]map[string]bool

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
func (m *Machine) canTransition(from string, to string) bool {
	return m.ruleSet[from][to]
}

// GetState get state
func (m *Machine) GetState(id string) (string, error) {
	state := &State{}
	if err := m.store.Get(context.Background(), path.Join(statesRootKey, id), state, true); err != nil {
		return "", err
	}
	return state.Value, nil
}

// TransitionTo transitionTo
func (m *Machine) TransitionTo(id string, to string) error {
	current, err := m.GetState(id)
	if err != nil {
		return err
	}
	if !m.canTransition(current, to) {
		return fmt.Errorf("Cannot transition from state %s to state %s", current, to)
	}
	expect := &State{Value: current}
	update := &State{Value: to}
	if err = m.store.CompareAndSet(context.Background(), path.Join(statesRootKey, id), expect, update); err != nil {
		return fmt.Errorf("Cannot transition from state %s to state %s", current, to)
	}
	return nil
}

// Is is
func (m *Machine) Is(id string, expected string) (bool, error) {
	state, err := m.GetState(id)
	if err != nil {
		return false, err
	}
	return state == expected, nil
}

// CreateState createstate
func (m *Machine) CreateState(id string, initial string) error {
	state := &State{Value: initial}
	if err := m.store.Create(context.Background(), path.Join(statesRootKey, id), state, nil, 0); err != nil {
		return err
	}
	return nil
}

// DeleteState deleteState
func (m *Machine) DeleteState(id string) error {
	return m.store.Delete(context.Background(), path.Join(statesRootKey, id), false, nil)
}
