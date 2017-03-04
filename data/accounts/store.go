package accounts

import (
	"context"
	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"github.com/hlandau/passlib"
	"github.com/ory-am/ladon"
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

func (s *Store) rawUser(ctx context.Context, name string) (*User, error) {
	user := &User{}
	if err := s.Store.Get(ctx, path.Join(usersRootKey, name), user, true); err != nil {
		return nil, err
	}
	// If there's no "name" in the answer, it means the user has not been found, so return nil
	if user.GetName() == "" {
		return nil, nil
	}
	return user, nil
}

func secureUser(user *User) *User {
	if user == nil {
		return nil
	}
	// For security reasons, remove the password hash
	user.PasswordHash = ""
	return user
}

func (s *Store) getUser(ctx context.Context, name string) (user *User, err error) {
	if user, err = s.rawUser(ctx, name); err != nil {
		return nil, err
	}
	if user == nil {
		return nil, UserNotFound
	}
	return user, nil
}

func (s *Store) getVerifiedUser(ctx context.Context, name string) (user *User, err error) {
	if user, err = s.getUser(ctx, name); err != nil {
		return nil, err
	}
	if !user.IsVerified {
		return nil, UserNotVerified
	}
	return user, nil
}

// Users

// CreateUser creates a new user
func (s *Store) CreateUser(ctx context.Context, name string, email string, password string) (user *User, err error) {
	// Check if user already exists
	userAlreadyExists, err := s.rawUser(ctx, name)
	if err != nil {
		return nil, err
	}
	if userAlreadyExists != nil {
		return nil, UserAlreadyExists
	}

	// Check if organization with the same name already exists
	orgAlreadyExists, err := s.GetOrganization(ctx, name)
	if err != nil {
		return nil, err
	}
	if orgAlreadyExists != nil {
		return nil, OrganizationAlreadyExists
	}

	// Create the new user
	user = &User{
		Email:      email,
		Name:       name,
		IsVerified: false,
		CreateDt:   time.Now().Unix(),
	}
	if err = CheckPassword(password); err != nil {
		return nil, err
	}
	if user.PasswordHash, err = passlib.Hash(password); err != nil {
		return nil, err
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := s.Store.Create(ctx, path.Join(usersRootKey, name), user, nil, 0); err != nil {
		return nil, err
	}
	return secureUser(user), nil
}

// VerifyUser verifies a user account
func (s *Store) VerifyUser(ctx context.Context, token string) (*User, error) {
	// Validate the token
	claims, err := auth.ValidateToken(token, auth.TokenTypeVerification)
	if err != nil {
		return nil, InvalidToken
	}
	user, err := s.getUser(ctx, claims.AccountName)
	if err != nil {
		return nil, err
	}
	user.IsVerified = true
	if err := s.Store.Put(ctx, path.Join(usersRootKey, user.Name), user, 0); err != nil {
		return nil, err
	}
	return secureUser(user), nil
}

// CheckUserPassword checks the given user password
func (s *Store) CheckUserPassword(ctx context.Context, name string, password string) error {
	user, err := s.getVerifiedUser(ctx, name)
	if err != nil {
		return err
	}
	// TODO: should we use the newHash ?
	_, err = passlib.Verify(password, user.PasswordHash)
	if err != nil {
		return WrongPassword
	}
	return nil
}

// SetUserPassword sets the given user password
func (s *Store) SetUserPassword(ctx context.Context, name string, password string) error {
	user, err := s.getUser(ctx, name)
	if err != nil {
		return err
	}

	// Password
	if err = CheckPassword(password); err != nil {
		return err
	}
	if user.PasswordHash, err = passlib.Hash(password); err != nil {
		return err
	}

	// Update user
	if err := s.Store.Put(ctx, path.Join(usersRootKey, user.Name), user, 0); err != nil {
		return err
	}
	return nil
}

// GetUser fetches a user by name
func (s *Store) GetUser(ctx context.Context, name string) (*User, error) {
	if err := CheckName(name); err != nil {
		return nil, err
	}
	user, err := s.rawUser(ctx, name)
	if err != nil {
		return nil, err
	}
	return secureUser(user), nil
}

// GetUserByEmail fetches a user by email
func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if _, err := CheckEmailAddress(email); err != nil {
		return nil, err
	}
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
func (s *Store) ListUsers(ctx context.Context) ([]*User, error) {
	protos := []proto.Message{}
	if err := s.Store.List(ctx, usersRootKey, storage.Everything, &User{}, &protos); err != nil {
		return nil, err
	}
	users := []*User{}
	for _, proto := range protos {
		users = append(users, secureUser(proto.(*User)))
	}
	return users, nil
}

// DeleteUser deletes a user by name
func (s *Store) DeleteUser(ctx context.Context, name string) error {
	// Get requester
	requester := auth.GetUser(ctx)
	if requester != name {
		return NotAuthorized
	}

	// Get organizations owned by he user
	ownedOrganizations, err := s.getOwnedOrganization(ctx, name)
	if err != nil {
		return err
	}
	// Check if user can be removed from all organizations
	for _, o := range ownedOrganizations {
		if _, err := s.canRemoveUserFromOrganization(ctx, o.Name, name); err != nil {
			return err
		}
	}
	// If yes, remove the user from all organizations
	for _, o := range ownedOrganizations {
		if err := s.RemoveUserFromOrganization(ctx, o.Name, name); err != nil {
			return err
		}
	}

	// Delete the user
	if err := s.Store.Delete(ctx, path.Join(usersRootKey, name), false, nil); err != nil {
		return err
	}
	return nil
}

// Organizations

func (s *Store) getOwnedOrganization(ctx context.Context, name string) ([]*Organization, error) {
	organizations, err := s.ListOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	ownedOrganizations := []*Organization{}
	for _, o := range organizations {
		if o.IsOwner(name) {
			ownedOrganizations = append(ownedOrganizations, o)
		}
	}
	return ownedOrganizations, nil
}

func (s *Store) updateOrganization(ctx context.Context, in *Organization) error {
	if err := in.Validate(); err != nil {
		return err
	}
	if err := s.Store.Put(ctx, path.Join(organizationsRootKey, in.Name), in, 0); err != nil {
		return err
	}
	return nil
}

// CreateOrganization creates a new organization
func (s *Store) CreateOrganization(ctx context.Context, name string, email string) error {
	// Check if user with the same name already exists
	userAlreadyExists, err := s.rawUser(ctx, name)
	if err != nil {
		return err
	}
	if userAlreadyExists != nil {
		return UserAlreadyExists
	}

	// Check if organization already exists
	orgAlreadyExists, err := s.GetOrganization(ctx, name)
	if err != nil {
		return err
	}
	if orgAlreadyExists != nil {
		return OrganizationAlreadyExists
	}

	// Create the new organization
	organization := &Organization{
		Email:    email,
		Name:     name,
		CreateDt: time.Now().Unix(),
		Members: []*OrganizationMember{
			{
				Name: auth.GetUser(ctx),
				Role: OrganizationRole_ORGANIZATION_OWNER,
			},
		},
	}
	if err := organization.Validate(); err != nil {
		return err
	}
	if err := s.Store.Create(ctx, path.Join(organizationsRootKey, organization.Name), organization, nil, 0); err != nil {
		return err
	}
	return nil
}

// AddUserToOrganization adds a user to the given organization
func (s *Store) AddUserToOrganization(ctx context.Context, organizationName string, userName string) (err error) {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return err
	}
	if organization == nil {
		return OrganizationNotFound
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  auth.GetUser(ctx),
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"resource": organization,
		},
	}); err != nil {
		return NotAuthorized
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Check if user is already a member
	if organization.HasMember(user.Name) {
		return nil // User is already a member of the organization, return
	}

	// Add the user as a team member
	organization.Members = append(organization.Members, &OrganizationMember{
		Name: user.Name,
		Role: OrganizationRole_ORGANIZATION_MEMBER,
	})
	return s.updateOrganization(ctx, organization)
}

func (s *Store) canRemoveUserFromOrganization(ctx context.Context, organizationName string, userName string) (*Organization, error) {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, OrganizationNotFound
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  auth.GetUser(ctx),
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"roles.organization": organization,
		},
	}); err != nil {
		return nil, NotAuthorized
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return nil, err
	}

	// Check if user is part of the organization
	memberIndex := organization.GetMemberIndex(user.Name)
	if memberIndex == -1 {
		return nil, nil // User is not a member of the organization, return
	}

	// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	organization.Members = append(organization.Members[:memberIndex], organization.Members[memberIndex+1:]...)
	if err := organization.Validate(); err != nil {
		return nil, err
	}
	return organization, nil
}

// RemoveUserFromOrganization removes a user from the given organization
func (s *Store) RemoveUserFromOrganization(ctx context.Context, organizationName string, userName string) (err error) {
	organization, err := s.canRemoveUserFromOrganization(ctx, organizationName, userName)
	if err != nil {
		return err
	}
	if organization == nil {
		return nil
	}
	return s.updateOrganization(ctx, organization)
}

// ChangeOrganizationMemberRole changes the role of given user in the given organization
func (s *Store) ChangeOrganizationMemberRole(ctx context.Context, organizationName string, userName string, role OrganizationRole) (err error) {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return err
	}
	if organization == nil {
		return OrganizationNotFound
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  auth.GetUser(ctx),
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"roles.organization": organization,
		},
	}); err != nil {
		return NotAuthorized
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Check if user is already a member
	member := organization.GetMember(user.Name)
	if member == nil {
		return UserNotFound
	}

	// Change the role of the user
	member.Role = role
	return s.updateOrganization(ctx, organization)
}

// GetOrganization fetches a organization by name
func (s *Store) GetOrganization(ctx context.Context, name string) (*Organization, error) {
	if err := CheckName(name); err != nil {
		return nil, err
	}
	organization := &Organization{}
	if err := s.Store.Get(ctx, path.Join(organizationsRootKey, name), organization, true); err != nil {
		return nil, err
	}
	// If there's no "name" in the answer, it means the organization has not been found, so return nil
	if organization.GetName() == "" {
		return nil, nil
	}
	return organization, nil
}

// ListOrganizations lists organizations
func (s *Store) ListOrganizations(ctx context.Context) ([]*Organization, error) {
	protos := []proto.Message{}
	if err := s.Store.List(ctx, organizationsRootKey, storage.Everything, &Organization{}, &protos); err != nil {
		return nil, err
	}
	organizations := []*Organization{}
	for _, proto := range protos {
		organizations = append(organizations, proto.(*Organization))
	}
	return organizations, nil
}

// DeleteOrganization deletes a organization by name
func (s *Store) DeleteOrganization(ctx context.Context, name string) error {
	// Get organization
	organization, err := s.GetOrganization(ctx, name)
	if err != nil {
		return err
	}
	if organization == nil {
		return OrganizationNotFound
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  auth.GetUser(ctx),
		Action:   auth.DeleteAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"resource": organization,
		},
	}); err != nil {
		return NotAuthorized
	}

	// Delete organization
	if err := s.Store.Delete(ctx, path.Join(organizationsRootKey, name), false, nil); err != nil {
		return err
	}
	return nil
}

// Teams

// CreateTeam creates a new team
func (s *Store) CreateTeam(ctx context.Context, organizationName, teamName string) error {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return err
	}
	if organization == nil {
		return OrganizationNotFound
	}

	// Check authorization
	requester := auth.GetUser(ctx)
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  requester,
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"resource": organization,
		},
	}); err != nil {
		return NotAuthorized
	}

	// Check if team already exists
	if organization.HasTeam(teamName) {
		return TeamAlreadyExists
	}

	// Create the new team
	team := &Team{
		Name:     teamName,
		CreateDt: time.Now().Unix(),
		Members: []*TeamMember{
			{
				Name: requester,
				Role: TeamRole_TEAM_OWNER,
			},
		},
	}

	// Add the team to the organization
	organization.Teams = append(organization.Teams, team)
	return s.updateOrganization(ctx, organization)
}

// AddUserToTeam adds a user to the given team
func (s *Store) AddUserToTeam(ctx context.Context, organizationName string, teamName string, userName string) error {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return err
	}
	if organization == nil {
		return OrganizationNotFound
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  auth.GetUser(ctx),
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"resource": organization,
		},
	}); err != nil {
		return NotAuthorized
	}

	// Get team
	team := organization.GetTeam(teamName)
	if team == nil {
		return TeamNotFound
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// TODO: Does the user need to be part of the organization?
	//// Check if user is part of the organization
	//if !organization.HasMember(user.Name) {
	//	return NotAnOrganizationMember
	//}

	// Check if user is part of the team
	if team.HasMember(user.Name) {
		return nil // User is already a member of the team, return
	}

	// Add the user as a team member
	team.Members = append(team.Members, &TeamMember{
		Name: user.Name,
		Role: TeamRole_TEAM_MEMBER,
	})
	return s.updateOrganization(ctx, organization)
}

// RemoveUserFromTeam removes a user from the given team
func (s *Store) RemoveUserFromTeam(ctx context.Context, organizationName string, teamName string, userName string) error {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return err
	}
	if organization == nil {
		return OrganizationNotFound
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  auth.GetUser(ctx),
		Action:   auth.UpdateAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"resource": organization,
		},
	}); err != nil {
		return NotAuthorized
	}

	// Get team
	team := organization.GetTeam(teamName)
	if team == nil {
		return TeamNotFound
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Check if user is actually a member
	memberIndex := team.GetMemberIndex(user.Name)
	if memberIndex == -1 {
		return nil // User is not a member of the team, return
	}

	// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	team.Members = append(team.Members[:memberIndex], team.Members[memberIndex+1:]...)
	return s.updateOrganization(ctx, organization)
}

// GetTeam fetches a team by name
func (s *Store) GetTeam(ctx context.Context, organizationName string, teamName string) (*Team, error) {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, OrganizationNotFound
	}

	// Get team
	if err := CheckName(teamName); err != nil {
		return nil, err
	}
	team := organization.GetTeam(teamName)
	return team, nil
}

// ListTeams lists teams
func (s *Store) ListTeams(ctx context.Context, organizationName string) ([]*Team, error) {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, OrganizationNotFound
	}
	return organization.Teams, nil
}

// DeleteTeam deletes a team by name
func (s *Store) DeleteTeam(ctx context.Context, organizationName string, teamName string) error {
	// Get organization
	organization, err := s.GetOrganization(ctx, organizationName)
	if err != nil {
		return err
	}
	if organization == nil {
		return OrganizationNotFound
	}

	// Check authorization
	if err := auth.Warden.IsAllowed(&ladon.Request{
		Subject:  auth.GetUser(ctx),
		Action:   auth.DeleteAction,
		Resource: auth.OrganizationResource,
		Context: ladon.Context{
			"resource": organization,
		},
	}); err != nil {
		return NotAuthorized
	}

	// Check if the team is actually a team in the organization
	teamIndex := organization.GetTeamIndex(teamName)
	if teamIndex == -1 {
		return nil // Team is not part of the organization, return
	}

	// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	organization.Teams = append(organization.Teams[:teamIndex], organization.Teams[teamIndex+1:]...)
	return s.updateOrganization(ctx, organization)
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) {
	s.Store.Delete(ctx, usersRootKey, true, nil)
	s.Store.Delete(ctx, organizationsRootKey, true, nil)
}
