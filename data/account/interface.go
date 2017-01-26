package account

import "github.com/appcelerator/amp/data/schema"

// Interface must be implemented an account database
type Interface interface {
	// AddAccount adds a new account to the account table
	AddAccount(account *schema.Account) (id string, err error)

	// Verify sets an account verification to true
	Verify(name string) error

	// AddTeam adds a new team to the team table
	AddTeam(team *schema.Team) (id string, err error)

	// AddTeamMember adds a new team to the team table
	AddTeamMember(teamId string, memberId string) (id string, err error)

	// GetAccount returns an account from the accounts table
	GetAccount(name string) (*schema.Account, error)

	// GetAccounts returns accounts matching a query
	GetAccounts(accountType schema.AccountType) ([]*schema.Account, error)
}
