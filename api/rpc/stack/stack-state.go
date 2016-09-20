package stack

import (
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/api/state"
)

// StackRuleSet defines possible transitions for stack states
var StackRuleSet = state.RuleSet{
	//   | Stopped   | Starting  | Running   | Redeploying
	[]bool{false /**/, true /* */, false /**/, true /* */}, // Stopped (initial state)
	[]bool{false /**/, false /**/, true /* */, false /**/}, // Starting
	[]bool{true /* */, false /**/, false /**/, true /* */}, // Running
	[]bool{true /* */, true /* */, false /**/, false /**/}, // Redeploying
}

var stackStateMachine = state.NewMachine(StackRuleSet, runtime.Store)
