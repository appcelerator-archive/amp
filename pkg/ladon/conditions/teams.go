package conditions

import (
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/ory-am/ladon"
)

// OwnerCondition is a condition which is fulfilled if the request's subject is an owner of the given resource
type TeamOwnerCondition struct{}

// Fulfills returns true if the request's subject is equal to the given value string
func (c *TeamOwnerCondition) Fulfills(value interface{}, r *ladon.Request) bool {
	members, ok := value.([]*schema.TeamMember)
	if !ok {
		return false
	}
	for _, member := range members {
		if member.Name == r.Subject {
			return true
		}
	}
	return false
}

// GetName returns the condition's name.
func (c *TeamOwnerCondition) GetName() string {
	return "TeamOwnerCondition"
}
