package stack

// RuleSet describe allowed state transitions
type RuleSet [][]bool

// StateMachine is the state machine
type StateMachine struct {
	ruleSet RuleSet
}

// NewStateMachine return a new state machine
func NewStateMachine(ruleSet RuleSet) StateMachine { return StateMachine{ruleSet: ruleSet} }

// CanTransition return whether or not you can transition between states
func (s *StateMachine) CanTransition(from int32, to int32) bool { return s.ruleSet[from][to] }
