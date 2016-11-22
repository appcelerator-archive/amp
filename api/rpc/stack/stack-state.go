package stack

import "github.com/appcelerator/amp/api/state"

// StackRuleSet defines possible transitions for stack states
var StackRuleSet = state.RuleSet{
	StackState_Stopped.String(): {
		StackState_Stopped.String():     false,
		StackState_Starting.String():    true,
		StackState_Running.String():     false,
		StackState_Redeploying.String(): true,
	},
	StackState_Starting.String(): {
		StackState_Stopped.String():     false,
		StackState_Starting.String():    false,
		StackState_Running.String():     true,
		StackState_Redeploying.String(): false,
	},
	StackState_Running.String(): {
		StackState_Stopped.String():     true,
		StackState_Starting.String():    false,
		StackState_Running.String():     false,
		StackState_Redeploying.String(): true,
	},
	StackState_Redeploying.String(): {
		StackState_Stopped.String():     true,
		StackState_Starting.String():    true,
		StackState_Running.String():     false,
		StackState_Redeploying.String(): false,
	},
}
