package auth

import (
	"github.com/appcelerator/amp/pkg/ladon/conditions"
	"github.com/docker/docker/pkg/stringid"
	"github.com/ory-am/ladon"
	"log"
)

const (
	// Resources
	AmpResource          = "amprn"
	OrganizationResource = AmpResource + ":organization"
	FunctionResource     = AmpResource + ":function"

	// Actions
	CreateAction = "create"
	ReadAction   = "read"
	UpdateAction = "update"
	DeleteAction = "delete"
	AdminAction  = CreateAction + "|" + ReadAction + "|" + UpdateAction + "|" + DeleteAction
)

var (
	// Organization owners are able to administrate their own organizations
	organizationOwners = &ladon.DefaultPolicy{
		ID:          stringid.GenerateNonCryptoID(),
		Subjects:    []string{"<.*>"},
		Description: "Organization owners are able to administrate their own organizations",
		Resources:   []string{OrganizationResource},
		Actions:     []string{"<" + AdminAction + ">"},
		Effect:      ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"resource": &conditions.ResourceOwnerCondition{},
		},
	}

	// Functions owners are able to administrate their own functions
	functionOwners = &ladon.DefaultPolicy{
		ID:          stringid.GenerateNonCryptoID(),
		Subjects:    []string{"<.*>"},
		Description: "Functions owners are able to administrate their own functions",
		Resources:   []string{FunctionResource},
		Actions:     []string{"<" + AdminAction + ">"},
		Effect:      ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"resource": &conditions.ResourceOwnerCondition{},
		},
	}

	// Policies represent access control policies for amp
	policies = []ladon.Policy{
		organizationOwners,
		functionOwners,
	}

	Warden = &ladon.Ladon{
		Manager: ladon.NewMemoryManager(),
	}
)

// TODO: Create a real policy manager?
func init() {
	// Register all policies
	for _, policy := range policies {
		if err := Warden.Manager.Create(policy); err != nil {
			log.Fatal("Unable to create policy:", err)
		}
	}
}
