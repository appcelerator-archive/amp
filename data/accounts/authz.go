package accounts

import (
	"log"

	"github.com/appcelerator/amp/api/auth"
	"github.com/docker/docker/pkg/stringid"
	"github.com/ory-am/ladon"
	"golang.org/x/net/context"
)

// Resources and actions
const (
	AmpResourceName = "amprn"
	UserRN          = AmpResourceName + ":user"
	OrganizationRN  = AmpResourceName + ":organization"
	TeamRN          = AmpResourceName + ":team"
	StackRN         = AmpResourceName + ":stack"

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
		Resources: []string{OrganizationRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"organization": &OrganizationAccessCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER},
				[]TeamPermissionLevel{},
			},
		},
	}

	teamsAdminByOrgOwners = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{TeamRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"organization": &OrganizationAccessCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER},
				[]TeamPermissionLevel{},
			},
		},
	}

	stacksAdminByOrgOwnersAndTeamAdmins = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{StackRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"organization": &OrganizationAccessCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER},
				[]TeamPermissionLevel{TeamPermissionLevel_TEAM_ADMIN},
			},
		},
	}

	stacksAdminByUserOwner = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{StackRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &ladon.EqualsSubjectCondition{},
		},
	}

	usersAdminByThemselves = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{UserRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &ladon.EqualsSubjectCondition{},
		},
	}

	// Policies represent access control policies for amp
	policies = []ladon.Policy{
		organizationsAdminByOrgOwners,
		stacksAdminByOrgOwnersAndTeamAdmins,
		stacksAdminByUserOwner,
		teamsAdminByOrgOwners,
		usersAdminByThemselves,
	}
)

// Authorization

// GetRequesterAccount gets the requester account from the given context, i.e. the user or organization performing the request
func GetRequesterAccount(ctx context.Context) *Account {
	activeOrganization := auth.GetActiveOrganization(ctx)
	if activeOrganization != "" {
		return &Account{
			Type: AccountType_ORGANIZATION,
			Name: activeOrganization,
		}
	}
	return &Account{
		Type: AccountType_USER,
		Name: auth.GetUser(ctx),
	}
}

// IsAuthorized returns whether the requesting user is authorized to perform the given action on given resource
func (s *Store) IsAuthorized(ctx context.Context, owner *Account, action string, resource string, resourceID string) bool {
	if owner == nil {
		return false
	}
	subject := auth.GetUser(ctx)
	switch owner.Type {
	case AccountType_ORGANIZATION:
		organization, err := s.getOrganization(ctx, owner.Name)
		if err != nil {
			return false
		}
		err = s.warden.IsAllowed(&ladon.Request{
			Subject:  subject,
			Action:   action,
			Resource: resource,
			Context: ladon.Context{
				"organization": organization,
				"resourceID":   resourceID,
			},
		})
		return err == nil
	case AccountType_USER:
		err := s.warden.IsAllowed(&ladon.Request{
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

// SuperOrganizationAccessCondition is a condition which is fulfilled if the request's subject is a member of the super organization
type SuperOrganizationAccessCondition struct {
	accounts *Store
}

// Fulfills returns true if subject is granted resource access
func (c *SuperOrganizationAccessCondition) Fulfills(value interface{}, r *ladon.Request) bool {
	so, err := c.accounts.GetOrganization(context.Background(), superOrganization)
	if err != nil {
		return false
	}
	return so.getMember(r.Subject) != nil
}

// GetName returns the condition's name.
func (c *SuperOrganizationAccessCondition) GetName() string {
	return "SuperOrganizationAccessCondition"
}

// OrganizationAccessCondition is a condition which is fulfilled if the request's subject has the expected access in the organization (either by organization role or team access)
type OrganizationAccessCondition struct {
	ExpectedRoles            []OrganizationRole
	ExpectedPermissionLevels []TeamPermissionLevel
}

// Fulfills returns true if subject is granted resource access
func (c *OrganizationAccessCondition) Fulfills(value interface{}, r *ladon.Request) bool {
	organization, ok := value.(*Organization)
	log.Println("organization", organization)
	if !ok {
		log.Println("organization ok:", ok)
		return false
	}
	if organization == nil {
		return false
	}
	member := organization.getMember(r.Subject)
	log.Println("member", member)
	if member == nil {
		return false
	}
	for _, role := range c.ExpectedRoles {
		if role == member.GetRole() {
			return true
		}
	}
	log.Println("organization.getMemberTeams(r.Subject)", organization.getMemberTeams(r.Subject))
	for _, team := range organization.getMemberTeams(r.Subject) {
		resourceID, ok := r.Context["resourceID"].(string)
		if !ok {
			log.Println("resourceID ok", ok)
			return false
		}
		log.Println("resourceID", resourceID)
		resource := team.getResourceById(resourceID)
		log.Println("resource", resource)
		if resource == nil {
			continue
		}
		log.Println("c.ExpectedPermissionLevels", c.ExpectedPermissionLevels)
		for _, pl := range c.ExpectedPermissionLevels {
			log.Println("pl", pl)
			log.Println("resource.GetPermissionLevel()", resource.GetPermissionLevel())
			log.Println("pl == resource.GetPermissionLevel()", pl == resource.GetPermissionLevel())
			if pl == resource.GetPermissionLevel() {
				return true
			}
		}
	}
	return false
}

// GetName returns the condition's name.
func (c *OrganizationAccessCondition) GetName() string {
	return "OrganizationRoleCondition"
}
