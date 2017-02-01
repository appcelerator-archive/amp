package account

import (
	"context"

	"github.com/appcelerator/amp/data/schema"
)

// Interface must be implemented an account database
type Interface interface {
	// AddAccount adds a new account to the account table
	AddAccount(ctx context.Context, account *schema.Account) (id string, err error)

	// Verify sets an account verification to true
	Verify(ctx context.Context, name string) error

	// AddTeam adds a new team to the team table
	AddTeam(ctx context.Context, team *schema.Team) (id string, err error)

	// AddTeamMember adds a new team to the team table
	AddTeamMember(ctx context.Context, teamID string, memberID string) (id string, err error)

	// GetAccount returns an account from the accounts table
	GetAccount(ctx context.Context, name string) (*schema.Account, error)

	// GetAccounts returns accounts matching a query
	GetAccounts(ctx context.Context, accountType schema.AccountType) ([]*schema.Account, error)
}
