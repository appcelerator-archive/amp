package account

import (
	"github.com/appcelerator/amp/data/storage"
	"golang.org/x/net/context"

	"github.com/appcelerator/amp/data/schema"
	"path"
)

// Mock impliments account data.Interface
type Etcd struct {
	Store storage.Interface
	ctx   context.Context
}

// NewMock returns a mock account database with some starter data
func NewEtcd(store storage.Interface, c context.Context) Interface {
	return &Etcd{
		Store: store,
		ctx:   c,
	}
}

// ETCD Implementation

// AddAccount adds a new account to the account table
func (e *Etcd) AddAccount(account *schema.Account) (id string, err error) {
	//TODO Add data integrity checks
	if err := e.Store.Create(e.ctx, path.Join(AccountRootNameKey, account.Id), account, nil, 0); err != nil {
		return "", err
	}
	return account.Id, nil

}

// Verify sets an account verification to true
func (e *Etcd) Verify(name string) error {
	acct, err := e.GetAccount(name)
	if err == nil && acct.Name != "" && !acct.IsVerified {
		acct.IsVerified = true
		err = e.Store.Put(e.ctx, path.Join(AccountRootNameKey, acct.Id), acct, 0)
	}
	return err
}

// AddTeam adds a new team to the team table
func (e *Etcd) AddTeam(team *schema.Team) (id string, err error) {
	return
}

// GetAccount returns an account from the accounts table
func (e *Etcd) GetAccount(name string) (account *schema.Account, err error) {
	acct := &schema.Account{}
	err = e.Store.Get(e.ctx, path.Join(AccountRootNameKey, name), acct, true)
	return acct, err
}

// GetAccounts implements Interface.GetAccounts
func (e *Etcd) GetAccounts(accountType schema.AccountType) (accounts []*schema.Account, err error) {
	return
}
