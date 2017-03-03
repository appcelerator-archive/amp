package account

import (
	"context"

	"github.com/appcelerator/amp/data/account/schema"
)

// Interface defines the user data access layer
type Interface interface {
	// CreateUser creates a new user with given password
	CreateUser(ctx context.Context, name string, email string, password string) (user *schema.User, err error)

	// CheckUserPassword checks the given user password
	CheckUserPassword(ctx context.Context, name string, password string) (err error)

	// SetUserPassword sets the given user password
	SetUserPassword(ctx context.Context, name string, password string) (err error)

	// GetUser fetches a user by name
	GetUser(ctx context.Context, name string) (user *schema.User, err error)

	// GetUserByEmail fetches a user by email
	GetUserByEmail(ctx context.Context, email string) (user *schema.User, err error)

	// ListUsers lists users
	ListUsers(ctx context.Context) (users []*schema.User, err error)

	// VerifyUser verifies a user account
	VerifyUser(ctx context.Context, token string) (user *schema.User, err error)

	// DeleteUser deletes the requester's user account
	DeleteUser(ctx context.Context) (user *schema.User, err error)

	// DeleteUser deletes a user by name
	DeleteUserByName(ctx context.Context, name string) (err error)

	// CreateOrganization creates a new organization
	CreateOrganization(ctx context.Context, name string, email string) (err error)

	// GetOrganization fetches a organization by name
	GetOrganization(ctx context.Context, name string) (organization *schema.Organization, err error)

	// AddUserToOrganization adds a user to the given organization
	AddUserToOrganization(ctx context.Context, organizationName string, userName string) (err error)

	// RemoveUserFromOrganization removes a user from the given organization
	RemoveUserFromOrganization(ctx context.Context, organizationName string, userName string) (err error)

	// ListOrganizations lists organizations
	ListOrganizations(ctx context.Context) (organizations []*schema.Organization, err error)

	// DeleteOrganization deletes a organization by name
	DeleteOrganization(ctx context.Context, name string) (err error)

	// CreateTeam creates a new team
	CreateTeam(ctx context.Context, organizationName string, teamName string) (err error)

	// GetTeam fetches a team by name
	GetTeam(ctx context.Context, organizationName string, teamName string) (team *schema.Team, err error)

	// ListTeams lists teams
	ListTeams(ctx context.Context, organizationName string) (teams []*schema.Team, err error)

	// AddUserToTeam adds a user to the given team
	AddUserToTeam(ctx context.Context, organizationName string, teamName string, userName string) (err error)

	// RemoveUserFromTeam removes a user from the given team
	RemoveUserFromTeam(ctx context.Context, organizationName string, teamName string, userName string) (err error)

	// DeleteTeam deletes a team by name
	DeleteTeam(ctx context.Context, organizationName string, teamName string) (err error)

	// Reset resets the user store
	Reset(ctx context.Context)
}
