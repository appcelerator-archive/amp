package account

import (
	"context"

	"github.com/appcelerator/amp/data/account/schema"
)

// Interface defines the user data access layer
type Interface interface {
	// CreateUser creates a new user with given password
	CreateUser(ctx context.Context, password string, user *schema.User) (err error)

	// CheckUserPassword checks the given user password
	CheckUserPassword(ctx context.Context, password string, name string) (err error)

	// SetUserPassword sets the given user password
	SetUserPassword(ctx context.Context, password string, name string) (err error)

	// GetUser fetches a user by name
	GetUser(ctx context.Context, name string) (user *schema.User, err error)

	// GetUserByEmail fetches a user by email
	GetUserByEmail(ctx context.Context, email string) (user *schema.User, err error)

	// ListUsers lists users
	ListUsers(ctx context.Context) (users []*schema.User, err error)

	// ActivateUser activates a user account
	ActivateUser(ctx context.Context, name string) (err error)

	// DeleteUser deletes a user by name
	DeleteUser(ctx context.Context, name string) (err error)

	// CreateOrganization creates a new organization
	CreateOrganization(ctx context.Context, organization *schema.Organization) (err error)

	// GetOrganization fetches a organization by name
	GetOrganization(ctx context.Context, name string) (organization *schema.Organization, err error)

	// GetOrganizationByEmail fetches a organization by email
	GetOrganizationByEmail(ctx context.Context, email string) (organization *schema.Organization, err error)

	// AddUserToOrganization adds a user to the given organization
	AddUserToOrganization(ctx context.Context, organization *schema.Organization, user *schema.User) (err error)

	// RemoveUserFromOrganization removes a user from the given organization
	RemoveUserFromOrganization(ctx context.Context, organization *schema.Organization, user *schema.User) (err error)

	// ListOrganizations lists organizations
	ListOrganizations(ctx context.Context) (organizations []*schema.Organization, err error)

	// DeleteOrganization deletes a organization by name
	DeleteOrganization(ctx context.Context, name string) (err error)

	// CreateTeam creates a new team
	CreateTeam(ctx context.Context, organization *schema.Organization, team *schema.Team) (err error)

	// AddUserToTeam adds a user to the given team
	AddUserToTeam(ctx context.Context, organization *schema.Organization, name string, user *schema.User) (err error)

	// RemoveUserFromTeam removes a user from the given team
	RemoveUserFromTeam(ctx context.Context, organization *schema.Organization, name string, user *schema.User) (err error)

	// DeleteTeam deletes a team by name
	DeleteTeam(ctx context.Context, organization *schema.Organization, name string) (err error)

	// Reset resets the user store
	Reset(ctx context.Context)
}
