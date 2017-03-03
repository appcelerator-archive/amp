package conditions

import (
	"github.com/ory-am/ladon"
)

type Resource interface {
	// GetOwners gets the resource owners
	GetOwners() (owners []string)
}

// ResourceOwnerCondition is a condition which is fulfilled if the request's subject is an owner of the given resource
type ResourceOwnerCondition struct{}

// Fulfills returns true if the request's subject is equal to the given value string
func (c *ResourceOwnerCondition) Fulfills(value interface{}, r *ladon.Request) bool {
	resource, ok := value.(Resource)
	if !ok {
		return false
	}
	for _, owner := range resource.GetOwners() {
		if owner == r.Subject {
			return true
		}
	}
	return false
}

// GetName returns the condition's name.
func (c *ResourceOwnerCondition) GetName() string {
	return "ResourceOwnerCondition"
}
