package account

import (
	"github.com/appcelerator/amp/data/storage"
	"golang.org/x/net/context"

	"fmt"
	"path"

	"github.com/appcelerator/amp/data/schema"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
)

// AccountSchemaRootKey the base key for all object within the accounts schema
const AccountSchemaRootKey = "accounts"

// AccountUserByNameKey stores the alternate key value by name
const AccountUserByNameKey = AccountSchemaRootKey + "/account/name"

// AccountUserByIdKey key used to store the account protobuf type
const AccountUserByIdKey = AccountSchemaRootKey + "/account/id"

//AccountTeamByNameKey key used to store the alternate key by name
const AccountTeamByNameKey = AccountSchemaRootKey + "/team/name"

//AccountTeamByIdKey key used to store the team protobuf type
const AccountTeamByIdKey = AccountSchemaRootKey + "/team/id"

// Store impliments account data.Interface
type Store struct {
	Store storage.Interface
	ctx   context.Context
}

// NewStore returns a Storage wrapper with functions to operate against the backing database
func NewStore(store storage.Interface, c context.Context) Interface {
	return &Store{
		Store: store,
		ctx:   c,
	}
}

//AddTeamMember adds a user account to the team table
func (s *Store) AddTeamMember(teamId string, memberId string) (id string, err error) {
	return
}

// generateUUID place holder until we standardize the approach we want to use
func generateUUID() (id string) {
	return stringid.GenerateNonCryptoID()
}

// AddAccount adds a new account to the account table
func (s *Store) AddAccount(account *schema.Account) (id string, err error) {

	if err = s.checkAccount(account); err == nil {
		// Store the Account struct and the alternate key
		if err = s.Store.Create(s.ctx, path.Join(AccountUserByIdKey, account.Id), account, nil, 0); err == nil {
			err = s.Store.Create(s.ctx, path.Join(AccountUserByNameKey, account.Name), &schema.ForeignKey{FkId: account.Id}, nil, 0)
		}
	}
	return account.Id, err
}
func (s *Store) checkAccount(account *schema.Account) error {

	acct, err := s.GetAccount(account.Name)
	if err == nil && acct.Id == "" {
		account.Id = generateUUID()

	} else {
		err = fmt.Errorf("Account %s already exists", acct.Name)
	}
	return err
}
func (s *Store) checkTeam(team *schema.Team) error {

	t, err := s.GetTeam(team.Name)
	if err == nil && t.Id == "" {
		team.Id = generateUUID()

	} else {
		err = fmt.Errorf("Team %s already exists", t.Name)
	}
	return err
}

// Verify sets an account verification to true
func (s *Store) Verify(name string) error {
	acct, err := s.GetAccount(name)
	if err == nil && acct.Name != "" && !acct.IsVerified {
		acct.IsVerified = true
		err = s.Store.Put(s.ctx, path.Join(AccountUserByIdKey, acct.Id), acct, 0)
	}
	return err
}

// AddTeam adds a new team to the team table
func (s *Store) AddTeam(team *schema.Team) (id string, err error) {
	//TODO Add data integrity checks
	if team.Id == "" {
		team.Id = generateUUID()
	}
	// Store Team struct and alternate Key
	if err = s.Store.Create(s.ctx, path.Join(AccountTeamByIdKey, team.Id), team, nil, 0); err == nil {
		err = s.Store.Create(s.ctx, path.Join(AccountTeamByNameKey, team.Name), &schema.ForeignKey{FkId: team.Id}, nil, 0)
	}
	return team.Id, err
}

// GetTeam returns a Team from the Team table
func (s *Store) GetTeam(name string) (team *schema.Team, err error) {
	team = &schema.Team{}
	fk := &schema.ForeignKey{}
	//Grab the ID
	err = s.Store.Get(s.ctx, path.Join(AccountTeamByNameKey, name), fk, true)
	if err == nil && fk.FkId != "" {
		err = s.Store.Get(s.ctx, path.Join(AccountTeamByIdKey, fk.FkId), team, true)
	}
	return team, err
}

// GetAccount returns an account from the accounts table
func (s *Store) GetAccount(name string) (account *schema.Account, err error) {
	acct := &schema.Account{}
	fk := &schema.ForeignKey{}
	//Grab the ID
	err = s.Store.Get(s.ctx, path.Join(AccountUserByNameKey, name), fk, true)
	if err == nil && fk.FkId != "" {
		err = s.Store.Get(s.ctx, path.Join(AccountUserByIdKey, fk.FkId), acct, true)
	}
	return acct, err
}

// GetAccounts implements Inrface.GetAccounts
func (s *Store) GetAccounts(accountType schema.AccountType) (accounts []*schema.Account, err error) {

	var out []proto.Message
	account := &schema.Account{}
	err = s.Store.List(s.ctx, AccountUserByIdKey, storage.Everything, account, &out)
	if err == nil {
		// Unfortunately we have to iterate and filter
		for i := 0; i < len(out); i++ {
			m, ok := out[i].(*schema.Account)
			if ok && m.Type == accountType {
				accounts = append(accounts, m)
			} else if !ok {
				err = fmt.Errorf("Unexpected Type Encountered")
				break
			}
		}
	}
	return
}
