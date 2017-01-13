package data

import (
	"strconv"

	"github.com/appcelerator/amp/data/schema"
)

// AddAccount adds a new account to the account table
func AddAccount(account *schema.Account) (id string, err error) {
	id = strconv.Itoa(len(mockAccounts))
	account.Id = id
	mockAccounts = append(mockAccounts, account)
	return
}

// Verify sets an account verification to true
func Verify(name string) error {
	for _, account := range mockAccounts {
		if account.Name == name {
			account.IsVerified = true
		}
	}
	return nil
}

// AddTeam adds a new team to the team table
func AddTeam(team *schema.Team) (id string, err error) {
	id = strconv.Itoa(len(mockTeams))
	team.Id = id
	mockTeams = append(mockTeams, team)
	return
}

// GetAccount returns an account from the accounts table
func GetAccount(name string) (account *schema.Account, err error) {
	for _, account := range mockAccounts {
		if account.Name == name {
			return account, nil
		}
	}
	return
}

//
func GetAccounts(accountType schema.AccountType) (accounts []*schema.Account, err error) {
	accounts = []*schema.Account{}
	for _, account := range accounts {
		if account.Type == accountType {
			accounts = append(accounts, account)
		}
	}
	return
}
