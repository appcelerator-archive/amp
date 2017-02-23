package account

import (
	"context"

	"github.com/appcelerator/amp/data/account/schema"
)

// Interface defines the user data access layer
type Interface interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *schema.User) (err error)

	// GetUser fetches a user by name
	GetUser(ctx context.Context, name string) (user *schema.User, err error)

	// GetUserByEmail fetches a user by email
	GetUserByEmail(ctx context.Context, email string) (user *schema.User, err error)

	// ListUsers lists users
	ListUsers(ctx context.Context) (users []*schema.User, err error)

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, update *schema.User) (err error)

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

	// Reset resets the user store
	Reset(ctx context.Context)
}
