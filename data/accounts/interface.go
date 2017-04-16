package accounts

import "golang.org/x/net/context"

// Error type
type Error string

func (e Error) Error() string {
	return string(e)
}

// Errors
const (
	InvalidName               = Error("username is invalid")
	InvalidEmail              = Error("email is invalid")
	PasswordTooWeak           = Error("password is too weak")
	WrongPassword             = Error("password is wrong")
	InvalidToken              = Error("token is invalid")
	UserAlreadyExists         = Error("user already exists")
	EmailAlreadyUsed          = Error("email is already in use")
	UserNotFound              = Error("user not found")
	UserNotVerified           = Error("user not verified")
	OrganizationAlreadyExists = Error("organization already exists")
	OrganizationNotFound      = Error("organization not found")
	TeamAlreadyExists         = Error("team already exists")
	TeamNotFound              = Error("team not found")
	AtLeastOneOwner           = Error("organization must have at least one owner")
	NotAuthorized             = Error("user not authorized")
	NotPartOfOrganization     = Error("user is not part of the organization")
	InvalidResourceID         = Error("invalid resource ID")
	ResourceNotFound          = Error("resource not found")
)

// Interface defines the user data access layer
type Interface interface {
	// CreateUser creates a new user with given password
	CreateUser(ctx context.Context, name string, email string, password string) (user *User, err error)

	// CheckUserPassword checks the given user password
	CheckUserPassword(ctx context.Context, name string, password string) (err error)

	// SetUserPassword sets the given user password
	SetUserPassword(ctx context.Context, name string, password string) (err error)

	// GetUser fetches a user by name
	GetUser(ctx context.Context, name string) (user *User, err error)

	// GetUserByEmail fetches a user by email
	GetUserByEmail(ctx context.Context, email string) (user *User, err error)

	// ListUsers lists users
	ListUsers(ctx context.Context) (users []*User, err error)

	// VerifyUser verifies a user account
	VerifyUser(ctx context.Context, token string) (user *User, err error)

	// DeleteUser deletes a user by name
	DeleteUser(ctx context.Context, name string) (err error)

	// CreateOrganization creates a new organization
	CreateOrganization(ctx context.Context, name string, email string) (err error)

	// GetOrganization fetches a organization by name
	GetOrganization(ctx context.Context, name string) (organization *Organization, err error)

	// AddUserToOrganization adds a user to the given organization
	AddUserToOrganization(ctx context.Context, organizationName string, userName string) (err error)

	// RemoveUserFromOrganization removes a user from the given organization
	RemoveUserFromOrganization(ctx context.Context, organizationName string, userName string) (err error)

	// ChangeOrganizationMemberRole changes the role of given user in the given organization
	ChangeOrganizationMemberRole(ctx context.Context, organizationName string, userName string, role OrganizationRole) (err error)

	// ListOrganizations lists organizations
	ListOrganizations(ctx context.Context) (organizations []*Organization, err error)

	// DeleteOrganization deletes a organization by name
	DeleteOrganization(ctx context.Context, name string) (err error)

	// CreateTeam creates a new team
	CreateTeam(ctx context.Context, organizationName string, teamName string) (err error)

	// GetTeam fetches a team by name
	GetTeam(ctx context.Context, organizationName string, teamName string) (team *Team, err error)

	// ListTeams lists teams
	ListTeams(ctx context.Context, organizationName string) (teams []*Team, err error)

	// AddUserToTeam adds a user to the given team
	AddUserToTeam(ctx context.Context, organizationName string, teamName string, userName string) (err error)

	// RemoveUserFromTeam removes a user from the given team
	RemoveUserFromTeam(ctx context.Context, organizationName string, teamName string, userName string) (err error)

	// AddResourceToTeam adds a resource to the given team
	AddResourceToTeam(ctx context.Context, organizationName string, teamName string, resourceName string) (err error)

	// RemoveResourceFromTeam removes a resource from the given team
	RemoveResourceFromTeam(ctx context.Context, organizationName string, teamName string, resourceName string) (err error)

	// ChangeTeamResourcePermissionLevel changes the permission level over the given resource in the given team
	ChangeTeamResourcePermissionLevel(ctx context.Context, organizationName string, teamName string, resource string, permissionLevel TeamPermissionLevel) (err error)

	// DeleteTeam deletes a team by name
	DeleteTeam(ctx context.Context, organizationName string, teamName string) (err error)

	// IsAuthorized returns whether the requesting user is authorized to perform the given action on given resource
	IsAuthorized(ctx context.Context, owner *Account, action string, resource string, resourceId string) bool

	// Reset resets the user store
	Reset(ctx context.Context)
}
