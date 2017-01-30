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

	//DeleteTeamMember
	DeleteTeamMember(teamId string, memberId string) (err error)

	// GetAccount returns an account from the accounts table
	GetAccount(name string) (*schema.Account, error)

	// GetAccounts returns accounts matching a query
	GetAccounts(accountType schema.AccountType) ([]*schema.Account, error)

	//AddResource Adds Resource to resource table
	AddResource(resource *schema.Resource) (id string, err error)

	//GetResource returns a team from the team table
	GetResource(name string) (team *schema.Resource, err error)

	//DeleteResource removes the Resource entry for a given Id
	DeleteResource(name string) (err error)

	//AddResource Adds Resource to resource table
	AddResourceSettings(resource *schema.ResourceSettings) (id string, err error)

	//GetResourceSettings returns a list of ResourceSettings for a given resource
	GetResourceSettings(resourceId string) (rs []*schema.ResourceSettings, err error)

	//DeleteResourceSettings removes the Resource entry for a given Id
	DeleteResourceSettings(resourceId string) (err error)

	//AddPermission Adds Permission to the Permission table
	AddPermission(resource *schema.Permission) (id string, err error)

	//GetPermission returns the permission record
	GetPermission(resourceId string) (perm []*schema.Permission, err error)
}
