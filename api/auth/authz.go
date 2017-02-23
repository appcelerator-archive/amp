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
	TeamResource         = AmpResource + ":team"

	// Actions
	CreateAction = "create"
	ReadAction   = "read"
	UpdateAction = "update"
	DeleteAction = "delete"
	AdminAction  = CreateAction + "|" + ReadAction + "|" + UpdateAction + "|" + DeleteAction
)

var (
	// Organization owners are able to administrate their own organizations
	organizationOwner = &ladon.DefaultPolicy{
		ID:          stringid.GenerateNonCryptoID(),
		Subjects:    []string{"<.*>"}, // This will match any subject (user name), we should consider using []string{"<.+>"} to have at least one character
		Description: "Organization owners are able to administrate their own organizations",
		Resources:   []string{OrganizationResource},
		Actions:     []string{"<" + AdminAction + ">"},
		Effect:      ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owners": &conditions.OrganizationOwnerCondition{},
		},
	}

	// Team owners are able to administrate their own teams
	teamOwner = &ladon.DefaultPolicy{
		ID:          stringid.GenerateNonCryptoID(),
		Subjects:    []string{"<.*>"}, // This will match any subject (user name), we should consider using []string{"<.+>"} to have at least one character
		Description: "Team owners are able to administrate their own teams",
		Resources:   []string{TeamResource},
		Actions:     []string{"<" + AdminAction + ">"},
		Effect:      ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owners": &conditions.TeamOwnerCondition{},
		},
	}

	// Policies represent access control policies for amp
	policies = []ladon.Policy{
		organizationOwner,
		teamOwner,
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
