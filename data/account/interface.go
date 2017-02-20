package account

import (
	"context"

	"github.com/appcelerator/amp/data/account/schema"
)

// Interface defines the user data access layer
type Interface interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *schema.User) (err error)

	// GetUser fetches an user by name
	GetUser(ctx context.Context, name string) (user *schema.User, err error)

	// GetUserByEmail fetches an user by email
	GetUserByEmail(ctx context.Context, email string) (user *schema.User, err error)

	// ListUsers lists users
	ListUsers(ctx context.Context) (users []*schema.User, err error)

	// UpdateUser updates an user
	UpdateUser(ctx context.Context, update *schema.User) (err error)

	// DeleteUser deletes an user by name
	DeleteUser(ctx context.Context, name string) (err error)

	// Reset resets the user store
	Reset(ctx context.Context) (err error)
}
