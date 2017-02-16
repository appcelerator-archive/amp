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

const usersRootKey = "users"
const organizationsRootKey = "organizations"

// Store implements user data.Interface
type Store struct {
	Store storage.Interface
}

// NewStore returns an etcd implementation of user.Interface
func NewStore(store storage.Interface) *Store {
	return &Store{
		Store: store,
	}
}

// CreateUser creates a new user
func (s *Store) CreateUser(ctx context.Context, in *schema.User) (string, error) {
	if err := in.Validate(); err != nil {
		return "", err
	}
	id := stringid.GenerateNonCryptoID()
	in.Id = id
	in.IsVerified = false
	in.CreateDt = time.Now().Unix()
	if err := s.Store.Create(ctx, path.Join(usersRootKey, id), in, nil, 0); err != nil {
		return "", err
	}
	return id, nil
}

// GetUser fetches an user by id
func (s *Store) GetUser(ctx context.Context, id string) (*schema.User, error) {
	user := &schema.User{}
	if err := s.Store.Get(ctx, path.Join(usersRootKey, id), user, false); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByName fetches an user by name
func (s *Store) GetUserByName(ctx context.Context, name string) (*schema.User, error) {
	users, err := s.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if strings.EqualFold(user.Name, name) {
			return user, nil
		}
	}
	return nil, nil
}

// GetUserByEmail fetches an user by email
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*schema.User, error) {
	users, err := s.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if strings.EqualFold(user.Email, email) {
			return user, nil
		}
	}
	return nil, nil
}

// ListUsers lists users
func (s *Store) ListUsers(ctx context.Context) ([]*schema.User, error) {
	protos := []proto.Message{}
	if err := s.Store.List(ctx, usersRootKey, storage.Everything, &schema.User{}, &protos); err != nil {
		return nil, err
	}
	users := []*schema.User{}
	for _, proto := range protos {
		users = append(users, proto.(*schema.User))
	}
	return users, nil
}

// UpdateUser updates an user
func (s *Store) UpdateUser(ctx context.Context, in *schema.User) error {
	if err := in.Validate(); err != nil {
		return err
	}
	if err := s.Store.Put(ctx, path.Join(usersRootKey, in.Id), in, 0); err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes an user by id
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	// TODO: check if user is owner of an organization
	if err := s.Store.Delete(ctx, path.Join(usersRootKey, id), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) error {
	if err := s.Store.Delete(ctx, path.Join(usersRootKey), true, nil); err != nil {
		return err
	}
	if err := s.Store.Delete(ctx, path.Join(organizationsRootKey), true, nil); err != nil {
		return err
	}
	return nil
}
