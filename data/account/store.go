package account

import (
	"context"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"path"
	"strings"
	"time"
)

const accountsRootKey = "accounts"

// Store implements account data.Interface
type Store struct {
	Store storage.Interface
}

// NewStore returns an etcd implementation of account.Interface
func NewStore(store storage.Interface) *Store {
	return &Store{
		Store: store,
	}
}

// CreateAccount creates a new account
func (s *Store) CreateAccount(ctx context.Context, in *schema.Account) (string, error) {
	id := stringid.GenerateNonCryptoID()
	in.Id = id
	in.CreateDt = time.Now().Unix()
	if err := s.Store.Create(ctx, path.Join(accountsRootKey, id), in, nil, 0); err != nil {
		return "", err
	}
	return id, nil
}

// GetAccount fetches an account by id
func (s *Store) GetAccount(ctx context.Context, id string) (*schema.Account, error) {
	account := &schema.Account{}
	if err := s.Store.Get(ctx, path.Join(accountsRootKey, id), account, false); err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountByUserName fetches an account by user name
func (s *Store) GetAccountByUserName(ctx context.Context, userName string) (*schema.Account, error) {
	accounts, err := s.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	for _, account := range accounts {
		if strings.EqualFold(account.UserName, userName) {
			return account, nil
		}
	}
	return nil, nil
}

// ListAccounts lists accounts
func (s *Store) ListAccounts(ctx context.Context) ([]*schema.Account, error) {
	protos := []proto.Message{}
	if err := s.Store.List(ctx, accountsRootKey, storage.Everything, &schema.Account{}, &protos); err != nil {
		return nil, err
	}
	accounts := []*schema.Account{}
	for _, proto := range protos {
		accounts = append(accounts, proto.(*schema.Account))
	}
	return accounts, nil
}

// UpdateAccount updates an account
func (s *Store) UpdateAccount(ctx context.Context, in *schema.Account) error {
	if err := s.Store.Put(ctx, path.Join(accountsRootKey, in.Id), in, 0); err != nil {
		return err
	}
	return nil
}

// DeleteAccount deletes an account by id
func (s *Store) DeleteAccount(ctx context.Context, id string) error {
	if err := s.Store.Delete(ctx, path.Join(accountsRootKey, id), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) error {
	if err := s.Store.Delete(ctx, path.Join(accountsRootKey), true, nil); err != nil {
		return err
	}
	return nil
}
