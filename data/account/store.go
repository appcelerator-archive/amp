package account

import (
	"context"
	"fmt"
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"
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
func (s *Store) CreateUser(ctx context.Context, in *schema.User) error {
	if err := in.Validate(); err != nil {
		return err
	}
	in.IsVerified = false
	in.CreateDt = time.Now().Unix()
	if err := s.Store.Create(ctx, path.Join(usersRootKey, in.Name), in, nil, 0); err != nil {
		return err
	}
	return nil
}

// GetUser fetches a user by name
func (s *Store) GetUser(ctx context.Context, name string) (*schema.User, error) {
	user := &schema.User{}
	if err := s.Store.Get(ctx, path.Join(usersRootKey, name), user, true); err != nil {
		return nil, err
	}
	// If there's no "name" in the answer, it means the user has not been found, so return nil
	if user.GetName() == "" {
		return nil, nil
	}
	return user, nil
}

// GetUserByEmail fetches a user by email
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

// GetUserFromContext fetches a user from context metadata
func (s *Store) GetUserFromContext(ctx context.Context) (*schema.User, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unable to get metadata from context")
	}
	users := md[auth.RequesterKey]
	if len(users) == 0 {
		return nil, fmt.Errorf("context metadata has no requester field")
	}
	user := users[0]
	return s.GetUser(ctx, user)
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

// UpdateUser updates a user
func (s *Store) UpdateUser(ctx context.Context, in *schema.User) error {
	if err := in.Validate(); err != nil {
		return err
	}
	if err := s.Store.Put(ctx, path.Join(usersRootKey, in.Name), in, 0); err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by name
func (s *Store) DeleteUser(ctx context.Context, name string) error {
	// TODO: check if user is owner of an organization
	if err := s.Store.Delete(ctx, path.Join(usersRootKey, name), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) error {
	if err := s.Store.Delete(ctx, usersRootKey, true, nil); err != nil {
		return err
	}
	if err := s.Store.Delete(ctx, organizationsRootKey, true, nil); err != nil {
		return err
	}
	return nil
}
