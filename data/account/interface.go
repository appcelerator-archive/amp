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

	//GetTeam returns a team from the team table
	GetTeam(name string) (team *schema.Team, err error)

	// AddTeamMember adds a new team to the team table
	AddTeamMember(teamMember *schema.TeamMember) (id string, err error)

	// GetTeamMember returns the TeamMember from the team_member table
	GetTeamMember(teamId string, memberId string) (member *schema.TeamMember, err error)

	// GetAccount returns an account from the accounts table
	GetAccount(name string) (*schema.Account, error)

	// GetAccounts returns accounts matching a query
	GetAccounts(accountType schema.AccountType) ([]*schema.Account, error)

	//AddResource Adds Resource to resource table
	AddResource(resource *schema.Resource) (id string, err error)

	//GetResource returns a team from the team table
	GetResource(name string) (team *schema.Resource, err error)

	//AddResource Adds Resource to resource table
	AddResourceSettings(resource *schema.ResourceSettings) (id string, err error)

	//GetResource returns a team from the team table
	//GetResource(name string) (team *schema.Resource, err error)

}
