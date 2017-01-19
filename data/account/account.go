package account

import (
	"strconv"

	"github.com/appcelerator/amp/data/schema"
)

const AccountRootNameKey = "accounts"

// AddAccount adds a new account to the account table
func (m *Mock) AddAccount(account *schema.Account) (id string, err error) {
	id = strconv.Itoa(len(m.accounts))
	account.Id = id
	m.accounts = append(m.accounts, account)
	return
}

// Verify sets an account verification to true
func (m *Mock) Verify(name string) error {
	for _, account := range m.accounts {
		if account.Name == name {
			account.IsVerified = true
		}
	}
	return nil
}

// AddTeam adds a new team to the team table
func (m *Mock) AddTeam(team *schema.Team) (id string, err error) {
	id = strconv.Itoa(len(m.teams))
	team.Id = id
	m.teams = append(m.teams, team)
	return
}

// GetAccount returns an account from the accounts table
func (m *Mock) GetAccount(name string) (account *schema.Account, err error) {
	for _, account := range m.accounts {
		if account.Name == name {
			return account, nil
		}
	}
	return
}

// GetAccounts implements Interface.GetAccounts
func (m *Mock) GetAccounts(accountType schema.AccountType) (accounts []*schema.Account, err error) {
	accounts = []*schema.Account{}
	for _, account := range m.accounts {
		if account.Type == accountType {
			accounts = append(accounts, account)
		}
	}
	return
}
