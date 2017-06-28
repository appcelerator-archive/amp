package accounts

import (
	log "github.com/Sirupsen/logrus"

	"github.com/appcelerator/amp/api/auth"
	"github.com/docker/docker/pkg/stringid"
	"github.com/ory/ladon"
	"golang.org/x/net/context"
)

// Resources and actions
const (
	AmpResourceName = "amprn"
	UserRN          = AmpResourceName + ":user"
	OrganizationRN  = AmpResourceName + ":organization"
	TeamRN          = AmpResourceName + ":team"
	StackRN         = AmpResourceName + ":stack"
	DashboardRN     = AmpResourceName + ":dashboard"

	CreateAction = "create"
	ReadAction   = "read"
	UpdateAction = "update"
	DeleteAction = "delete"
	LeaveAction  = "leave"
	AdminAction  = "admin"
	AnyAction    = CreateAction + "|" + ReadAction + "|" + UpdateAction + "|" + DeleteAction + "|" + LeaveAction + "|" + AdminAction
)

var (
	usersAdminByThemselves = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{UserRN},
		Actions:   []string{"<" + AnyAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{},
		},
	}

	usersCanLeaveOrganizations = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{OrganizationRN},
		Actions:   []string{"<" + LeaveAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{},
		},
	}

	organizationsAdminByOrgOwners = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{OrganizationRN},
		Actions:   []string{"<" + AnyAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER},
				[]TeamPermissionLevel{},
			},
		},
	}

	teamsAdminByOrgOwners = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{TeamRN},
		Actions:   []string{"<" + AnyAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER},
				[]TeamPermissionLevel{},
			},
		},
	}

	teamsCanBeCreatedByOrgMembers = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{TeamRN},
		Actions:   []string{"<" + CreateAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER, OrganizationRole_ORGANIZATION_MEMBER},
				[]TeamPermissionLevel{},
			},
		},
	}

	stacksReadByTeamReadPermission = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{StackRN},
		Actions:   []string{"<" + ReadAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{},
				[]TeamPermissionLevel{TeamPermissionLevel_TEAM_READ},
			},
		},
	}

	stacksReadUpdateByTeamWritePermission = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{StackRN},
		Actions:   []string{"<" + ReadAction + "|" + UpdateAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{},
				[]TeamPermissionLevel{TeamPermissionLevel_TEAM_WRITE},
			},
		},
	}

	stacksAdminByOwner = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{StackRN},
		Actions:   []string{"<" + AnyAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER},
				[]TeamPermissionLevel{TeamPermissionLevel_TEAM_ADMIN},
			},
		},
	}

	dashboardsReadByTeamReadPermission = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{DashboardRN},
		Actions:   []string{"<" + ReadAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{},
				[]TeamPermissionLevel{TeamPermissionLevel_TEAM_READ},
			},
		},
	}

	dashboardsReadUpdateByTeamWritePermission = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{DashboardRN},
		Actions:   []string{"<" + ReadAction + "|" + UpdateAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{},
				[]TeamPermissionLevel{TeamPermissionLevel_TEAM_WRITE},
			},
		},
	}

	dashboardsAdminByOwners = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{DashboardRN},
		Actions:   []string{"<" + AnyAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &OwnerCondition{
				[]OrganizationRole{OrganizationRole_ORGANIZATION_OWNER},
				[]TeamPermissionLevel{TeamPermissionLevel_TEAM_ADMIN},
			},
		},
	}

	// Policies represent access control policies for amp
	policies = []ladon.Policy{
		usersAdminByThemselves,
		usersCanLeaveOrganizations,
		organizationsAdminByOrgOwners,
		teamsAdminByOrgOwners,
		teamsCanBeCreatedByOrgMembers,
		stacksReadByTeamReadPermission,
		stacksReadUpdateByTeamWritePermission,
		stacksAdminByOwner,
		dashboardsReadByTeamReadPermission,
		dashboardsReadUpdateByTeamWritePermission,
		dashboardsAdminByOwners,
	}
)

// Authorization

// GetRequesterAccount gets the requester account from the given context, i.e. the user or organization performing the request
func GetRequesterAccount(ctx context.Context) *Account {
	return &Account{
		User:         auth.GetUser(ctx),
		Organization: auth.GetActiveOrganization(ctx),
	}
}

// IsAuthorized returns whether the requesting user is authorized to perform the given action on given resource
func (s *Store) IsAuthorized(ctx context.Context, owner *Account, action string, resource string, resourceID string) bool {
	subject := auth.GetUser(ctx)
	log.Debugf("IsAuthorized: ctx(subject): %s, owner: %v, action: %s, resource: %s, resourceID: %s\n", subject, owner, action, resource, resourceID)
	if owner == nil {
		return false
	}
	err := s.warden.IsAllowed(&ladon.Request{
		Subject:  subject,
		Action:   action,
		Resource: resource,
		Context: ladon.Context{
			"owner":      owner,
			"store":      s,
			"resourceID": resourceID,
		},
	})
	return err == nil
}

// OwnerCondition is a condition which is fulfilled if the request's subject has ownership over the resource
type OwnerCondition struct {
	ExpectedRoles            []OrganizationRole
	ExpectedPermissionLevels []TeamPermissionLevel
}

// Fulfills returns true if subject is granted resource access
func (c *OwnerCondition) Fulfills(value interface{}, r *ladon.Request) bool {
	// Validate context
	owner, ok := r.Context["owner"].(*Account)
	if !ok {
		return false
	}
	store, ok := r.Context["store"].(*Store)
	if !ok {
		return false
	}
	resourceID, ok := r.Context["resourceID"].(string)
	if !ok {
		return false
	}

	// Get super organization
	so, err := store.GetOrganization(context.Background(), superOrganization)
	if err != nil {
		return false
	}

	// The subject has full access if it's member of the super organization
	if so.getMember(r.Subject) != nil {
		return true
	}

	// The subject has full access if it's the owner of the resource
	if owner.User == r.Subject {
		return true
	}

	// If the resource is not created inside an organization, just return
	if owner.Organization == "" {
		return false
	}

	// Retrieve the owning organization
	organization, err := store.getOrganization(context.Background(), owner.Organization)
	if err != nil {
		return false
	}
	if organization == nil {
		return false
	}

	// Make sure the subject is a member of the owning organization
	member := organization.getMember(r.Subject)
	if member == nil {
		return false
	}

	// The subject has access with an expected expectedRole
	for _, expectedRole := range c.ExpectedRoles {
		if expectedRole == member.GetRole() {
			return true
		}
	}

	// The subject is only a member of the organization, check team access permissions
	for _, team := range organization.getMemberTeams(r.Subject) {
		resource := team.getResourceById(resourceID)
		if resource == nil {
			continue
		}
		for _, pl := range c.ExpectedPermissionLevels {
			if pl == resource.GetPermissionLevel() {
				return true
			}
		}
	}
	return false
}

// GetName returns the condition's name.
func (c *OwnerCondition) GetName() string {
	return "OwnerCondition"
}
