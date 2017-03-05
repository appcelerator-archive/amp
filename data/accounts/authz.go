package accounts

import (
	"github.com/appcelerator/amp/api/authn"
	"github.com/docker/docker/pkg/stringid"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
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

// Authorization

// GetRequesterAccount gets the requester account from the given context, i.e. the user or organization performing the request
func GetRequesterAccount(ctx context.Context) *Account {
	activeOrganization := authn.GetActiveOrganization(ctx)
	if activeOrganization != "" {
		return &Account{
			Type: AccountType_ORGANIZATION,
			Name: activeOrganization,
		}
	}
	return &Account{
		Type: AccountType_USER,
		Name: authn.GetUser(ctx),
	}
}

// IsAuthorized returns whether the requesting user is authorized to perform the given action on given resource
func (s *Store) IsAuthorized(ctx context.Context, owner *Account, action string, resource string) bool {
	if owner == nil {
		return false
	}
	subject := authn.GetUser(ctx)
	switch owner.Type {
	case AccountType_ORGANIZATION:
		organization, err := s.GetOrganization(ctx, owner.Name)
		if err != nil {
			return false
		}
		if organization == nil {
			return false
		}
		err = warden.IsAllowed(&ladon.Request{
			Subject:  subject,
			Action:   action,
			Resource: resource,
			Context: ladon.Context{
				"organization": organization,
			},
		})
		return err == nil
	case AccountType_USER:
		err := warden.IsAllowed(&ladon.Request{
			Subject:  subject,
			Action:   action,
			Resource: resource,
			Context: ladon.Context{
				"user": owner.Name,
			},
		})
		return err == nil
	}
	return false
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
