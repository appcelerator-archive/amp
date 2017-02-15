package account

import (
	"context"

	"github.com/appcelerator/amp/data/account/schema"
)

// Interface defines the Account data access layer
type Interface interface {
	// CreateAccount creates a new account
	CreateAccount(ctx context.Context, account *schema.Account) (id string, err error)

	// GetAccount fetches an account by id
	GetAccount(ctx context.Context, id string) (account *schema.Account, err error)

	// GetAccountByUserName fetches an account by user name
	GetAccountByUserName(ctx context.Context, userName string) (account *schema.Account, err error)

	// ListAccounts lists accounts
	ListAccounts(ctx context.Context) (accounts []*schema.Account, err error)

	// UpdateAccount updates an account
	UpdateAccount(ctx context.Context, update *schema.Account) (err error)

	// DeleteAccount deletes an account by id
	DeleteAccount(ctx context.Context, id string) (err error)

	// Reset resets the account store
	Reset(ctx context.Context) (err error)
}
