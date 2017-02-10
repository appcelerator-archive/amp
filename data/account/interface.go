package account

import (
	"context"

	"github.com/appcelerator/amp/data/account/schema"
)

// Interface defines the user data access layer
type Interface interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *schema.User) (id string, err error)

	// GetUser fetches an user by id
	GetUser(ctx context.Context, id string) (user *schema.User, err error)

	// GetUserByName fetches an user by name
	GetUserByName(ctx context.Context, name string) (user *schema.User, err error)

	// GetUserByEmail fetches an user by email
	GetUserByEmail(ctx context.Context, email string) (user *schema.User, err error)

	// ListUsers lists users
	ListUsers(ctx context.Context) (users []*schema.User, err error)

	// UpdateUser updates an user
	UpdateUser(ctx context.Context, update *schema.User) (err error)

	// DeleteUser deletes an user by id
	DeleteUser(ctx context.Context, id string) (err error)

	// Reset resets the user store
	Reset(ctx context.Context) (err error)
}
