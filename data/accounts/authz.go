package accounts

import (
	"github.com/docker/docker/pkg/stringid"
	"github.com/ory-am/ladon"
	"log"
)

const (
	// Resources
	AmpResource          = "amprn"
	OrganizationResource = AmpResource + ":organization"
	TeamResource         = AmpResource + ":team"
	FunctionResource     = AmpResource + ":function"

	// Actions
	CreateAction = "create"
	ReadAction   = "read"
	UpdateAction = "update"
	DeleteAction = "delete"
	AdminAction  = CreateAction + "|" + ReadAction + "|" + UpdateAction + "|" + DeleteAction
)

var (
	organizationsAdminByOrgOwners = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{OrganizationResource},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"organization": &OrganizationRoleCondition{[]OrganizationRole{
				OrganizationRole_ORGANIZATION_OWNER,
			}},
		},
	}

	teamsAdminByOrgOwnersAndMembers = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{TeamResource},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"organization": &OrganizationRoleCondition{[]OrganizationRole{
				OrganizationRole_ORGANIZATION_MEMBER,
				OrganizationRole_ORGANIZATION_OWNER,
			}},
		},
	}

	functionsAdminByOrgOwnersAndMembers = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{FunctionResource},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"organization": &OrganizationRoleCondition{[]OrganizationRole{
				OrganizationRole_ORGANIZATION_MEMBER,
				OrganizationRole_ORGANIZATION_OWNER,
			}},
		},
	}

	functionsAdminByUserOwner = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{FunctionResource},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &ladon.EqualsSubjectCondition{},
		},
	}

	// Policies represent access control policies for amp
	policies = []ladon.Policy{
		organizationsAdminByOrgOwners,
		teamsAdminByOrgOwnersAndMembers,
		functionsAdminByOrgOwnersAndMembers,
		functionsAdminByUserOwner,
	}

	warden = &ladon.Ladon{
		Manager: ladon.NewMemoryManager(),
	}
)

// TODO: Create a real policy manager?
func init() {
	// Register all policies
	for _, policy := range policies {
		if err := warden.Manager.Create(policy); err != nil {
			log.Fatal("Unable to create policy:", err)
		}
	}
}

// OrganizationRoleCondition is a condition which is fulfilled if the request's subject has the expected role in the organization
type OrganizationRoleCondition struct {
	ExpectedRoles []OrganizationRole
}

// Fulfills returns true if the request's subject is equal to the given value string
func (c *OrganizationRoleCondition) Fulfills(value interface{}, r *ladon.Request) bool {
	organization, ok := value.(*Organization)
	if !ok {
		return false
	}
	if organization == nil {
		return false
	}
	member := organization.GetMember(r.Subject)
	if member == nil {
		return false
	}
	for _, role := range c.ExpectedRoles {
		if role == member.GetRole() {
			return true
		}
	}
	return false
}

// GetName returns the condition's name.
func (c *OrganizationRoleCondition) GetName() string {
	return "OrganizationRoleCondition"
}
