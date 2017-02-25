package account

import (
	"context"
	"fmt"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"github.com/hlandau/passlib"
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
func (s *Store) CreateUser(ctx context.Context, password string, in *schema.User) (err error) {
	if err := in.Validate(); err != nil {
		return err
	}
	in.IsVerified = false
	in.CreateDt = time.Now().Unix()
	in.PasswordHash, err = passlib.Hash(password)
	if err != nil {
		return err
	}
	if err := s.Store.Create(ctx, path.Join(usersRootKey, in.Name), in, nil, 0); err != nil {
		return err
	}
	return nil
}

func (s *Store) rawUser(ctx context.Context, name string) (*schema.User, error) {
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

func secureUser(user *schema.User) *schema.User {
	if user == nil {
		return nil
	}
	// For security reasons, remove the password hash
	user.PasswordHash = ""
	return user
}

// GetUser fetches a user by name
func (s *Store) GetUser(ctx context.Context, name string) (*schema.User, error) {
	user, err := s.rawUser(ctx, name)
	if err != nil {
		return nil, err
	}
	return secureUser(user), nil
}

// CheckUserPassword checks the given user password
func (s *Store) CheckUserPassword(ctx context.Context, password string, name string) error {
	user, err := s.rawUser(ctx, name)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	// TODO: should we use the newHash ?
	_, err = passlib.Verify(password, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}

// SetUserPassword sets the given user password
func (s *Store) SetUserPassword(ctx context.Context, password string, name string) error {
	user, err := s.rawUser(ctx, name)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	user.PasswordHash, err = passlib.Hash(password)
	if err != nil {
		return err
	}
	if err := s.Store.Put(ctx, path.Join(usersRootKey, user.Name), user, 0); err != nil {
		return err
	}
	return nil
}

// GetUserByEmail fetches a user by email
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*schema.User, error) {
	users, err := s.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if strings.EqualFold(user.Email, email) {
			return secureUser(user), nil
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
		users = append(users, secureUser(proto.(*schema.User)))
	}
	return users, nil
}

// ActivateUser activates a user account
func (s *Store) ActivateUser(ctx context.Context, name string) error {
	user, err := s.rawUser(ctx, name)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	user.IsVerified = true
	if err := s.Store.Put(ctx, path.Join(usersRootKey, user.Name), user, 0); err != nil {
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
		Role: schema.OrganizationRole_ORGANIZATION_MEMBER,
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

// CreateTeam creates a new team
func (s *Store) CreateTeam(ctx context.Context, organization *schema.Organization, team *schema.Team) error {
	if err := team.Validate(); err != nil {
		return err
	}
	// Check if team already exists
	for _, t := range organization.Teams {
		if t.Name == team.Name {
			return fmt.Errorf("team already exists")
		}
	}
	team.CreateDt = time.Now().Unix()
	// Add the team to the organization
	organization.Teams = append(organization.Teams, team)
	return s.updateOrganization(ctx, organization)
}

// AddUserToTeam adds a user to the given team
func (s *Store) AddUserToTeam(ctx context.Context, organization *schema.Organization, teamName string, user *schema.User) error {
	team := organization.GetTeam(teamName)
	if team == nil {
		return fmt.Errorf("team not found")
	}

	// Check if user is already a member
	for _, member := range team.Members {
		if member.Name == user.Name {
			return nil // User is already a member of the team, return
		}
	}

	// Add the user as a team member
	team.Members = append(team.Members, &schema.TeamMember{
		Name: user.Name,
		Role: schema.TeamRole_TEAM_MEMBER,
	})

	return s.updateOrganization(ctx, organization)
}

// RemoveUserFromTeam removes a user from the given team
func (s *Store) RemoveUserFromTeam(ctx context.Context, organization *schema.Organization, teamName string, user *schema.User) error {
	team := organization.GetTeam(teamName)
	if team == nil {
		return fmt.Errorf("team not found")
	}

	// Check if user is actually a member
	memberIndex := -1
	for i, member := range team.Members {
		if member.Name == user.Name {
			memberIndex = i
			break
		}
	}
	if memberIndex == -1 {
		return nil // User is not a member of the team, return
	}

	// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	team.Members = append(team.Members[:memberIndex], team.Members[memberIndex+1:]...)
	return s.updateOrganization(ctx, organization)
}

// DeleteTeam deletes a team by name
func (s *Store) DeleteTeam(ctx context.Context, organization *schema.Organization, teamName string) error {
	team := organization.GetTeam(teamName)
	if team == nil {
		return fmt.Errorf("team not found")
	}

	// Check if the team is actually a team in the organization
	teamIndex := -1
	for i, team := range team.Members {
		if team.Name == teamName {
			teamIndex = i
			break
		}
	}
	if teamIndex == -1 {
		return nil // Team is not part of the organization team, return
	}

	// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	organization.Members = append(organization.Members[:teamIndex], organization.Members[teamIndex+1:]...)
	return s.updateOrganization(ctx, organization)
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) {
	s.Store.Delete(ctx, usersRootKey, true, nil)
	s.Store.Delete(ctx, organizationsRootKey, true, nil)
}
