package account

import (
	"context"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/data/storage"
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

// Users

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

// Organizations

// CreateOrganization creates a new organization
func (s *Store) CreateOrganization(ctx context.Context, in *schema.Organization) error {
	if err := in.Validate(); err != nil {
		return err
	}
	in.CreateDt = time.Now().Unix()
	if err := s.Store.Create(ctx, path.Join(organizationsRootKey, in.Name), in, nil, 0); err != nil {
		return err
	}
	return nil
}

// GetOrganization fetches a organization by name
func (s *Store) GetOrganization(ctx context.Context, name string) (*schema.Organization, error) {
	organization := &schema.Organization{}
	if err := s.Store.Get(ctx, path.Join(organizationsRootKey, name), organization, true); err != nil {
		return nil, err
	}
	// If there's no "name" in the answer, it means the organization has not been found, so return nil
	if organization.GetName() == "" {
		return nil, nil
	}
	return organization, nil
}

// GetOrganizationByEmail fetches a organization by email
func (s *Store) GetOrganizationByEmail(ctx context.Context, email string) (*schema.Organization, error) {
	organizations, err := s.ListOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	for _, organization := range organizations {
		if strings.EqualFold(organization.Email, email) {
			return organization, nil
		}
	}
	return nil, nil
}

// AddUserToOrganization adds a user to the given organization
func (s *Store) AddUserToOrganization(ctx context.Context, organization *schema.Organization, user *schema.User) (err error) {
	// Check if user is already a member
	for _, member := range organization.Members {
		if member.Name == user.Name {
			return nil // User is already a member of the organization, return
		}
	}
	// Add the user as a team member
	organization.Members = append(organization.Members, &schema.OrganizationMember{
		Name: user.Name,
		Role: schema.OrganizationRole_MEMBER,
	})
	return s.updateOrganization(ctx, organization)
}

// RemoveUserFromOrganization removes a user from the given organization
func (s *Store) RemoveUserFromOrganization(ctx context.Context, organization *schema.Organization, user *schema.User) (err error) {
	// Check if user is actually a member
	memberIndex := -1
	for i, member := range organization.Members {
		if member.Name == user.Name {
			memberIndex = i
			break
		}
	}
	if memberIndex == -1 {
		return nil // User is not a member of the organization, return
	}

	// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	organization.Members = append(organization.Members[:memberIndex], organization.Members[memberIndex+1:]...)
	return s.updateOrganization(ctx, organization)
}

// ListOrganizations lists organizations
func (s *Store) ListOrganizations(ctx context.Context) ([]*schema.Organization, error) {
	protos := []proto.Message{}
	if err := s.Store.List(ctx, organizationsRootKey, storage.Everything, &schema.Organization{}, &protos); err != nil {
		return nil, err
	}
	organizations := []*schema.Organization{}
	for _, proto := range protos {
		organizations = append(organizations, proto.(*schema.Organization))
	}
	return organizations, nil
}

// updateOrganization updates a organization
func (s *Store) updateOrganization(ctx context.Context, in *schema.Organization) error {
	if err := in.Validate(); err != nil {
		return err
	}
	if err := s.Store.Put(ctx, path.Join(organizationsRootKey, in.Name), in, 0); err != nil {
		return err
	}
	return nil
}

// DeleteOrganization deletes a organization by name
func (s *Store) DeleteOrganization(ctx context.Context, name string) error {
	// TODO: check preconditions
	if err := s.Store.Delete(ctx, path.Join(organizationsRootKey, name), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) {
	s.Store.Delete(ctx, usersRootKey, true, nil)
	s.Store.Delete(ctx, organizationsRootKey, true, nil)
}
