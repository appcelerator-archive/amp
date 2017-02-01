package account

import (
	"context"

	"github.com/appcelerator/amp/data/storage"

	"fmt"
	"path"

	"github.com/appcelerator/amp/data/schema"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
)

const accountSchemaRootKey = "accounts"
const accountUserByNameKey = accountSchemaRootKey + "/account/name"
const accountUserByID = accountSchemaRootKey + "/account/id"
const accountTeamKey = accountSchemaRootKey + "/team"

// Store impliments account data.Interface
type Store struct {
	Store storage.Interface
}

// NewStore returns a Storage wrapper with functions to operate against the backing database
func NewStore(store storage.Interface) *Store {
	return &Store{
		Store: store,
	}
}

//AddTeamMember adds a user account to the team table
func (s *Store) AddTeamMember(ctx context.Context, teamID string, memberID string) (id string, err error) {
	return
}

// generateUUID place holder until we standardize the approach we want to use
func generateUUID() (id string) {
	return stringid.GenerateNonCryptoID()
}

// AddAccount adds a new account to the account table
func (s *Store) AddAccount(ctx context.Context, account *schema.Account) (id string, err error) {

	if err = s.checkAccount(ctx, account); err == nil {
		// Store the account struct and the alternate key
		if err = s.Store.Create(ctx, path.Join(accountUserByID, account.Id), account, nil, 0); err == nil {
			fk := &schema.ForeignKey{FkId: account.Id}
			err = s.Store.Create(ctx, path.Join(accountUserByNameKey, account.Name), fk, nil, 0)
		}
	}
	return account.Id, err
}
func (s *Store) checkAccount(ctx context.Context, account *schema.Account) error {

	acct, err := s.GetAccount(ctx, account.Name)
	if err == nil && acct.Id == "" {
		account.Id = generateUUID()
	} else {
		err = fmt.Errorf("Account %s already exists", acct.Name)
	}
	return err
}

// Verify sets an account verification to true
func (s *Store) Verify(ctx context.Context, name string) error {
	acct, err := s.GetAccount(ctx, name)
	if err == nil && acct.Name != "" && !acct.IsVerified {
		acct.IsVerified = true
		err = s.Store.Put(ctx, path.Join(accountUserByID, acct.Id), acct, 0)
	}
	return err
}

// AddTeam adds a new team to the team table
func (s *Store) AddTeam(ctx context.Context, team *schema.Team) (id string, err error) {
	//TODO Add data integrity checks
	if team.Id == "" {
		team.Id = generateUUID()
	}
	if err := s.Store.Create(ctx, path.Join(accountSchemaRootKey, team.Id), team, nil, 0); err != nil {
		return "", err
	}
	return team.Id, nil
}

// GetAccount returns an account from the accounts table
func (s *Store) GetAccount(ctx context.Context, name string) (account *schema.Account, err error) {
	acct := &schema.Account{}
	fk := &schema.ForeignKey{}
	//Grab the ID
	err = s.Store.Get(ctx, path.Join(accountUserByNameKey, name), fk, true)
	if err == nil && fk.FkId != "" {
		err = s.Store.Get(ctx, path.Join(accountUserByID, fk.FkId), acct, true)
	}
	return acct, err
}

// GetAccounts implements Inrface.GetAccounts
func (s *Store) GetAccounts(ctx context.Context, accountType schema.AccountType) (accounts []*schema.Account, err error) {

	var out []proto.Message
	account := &schema.Account{}
	err = s.Store.List(ctx, accountUserByID, storage.Everything, account, &out)
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
