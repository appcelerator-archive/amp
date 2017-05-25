package accounts

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"github.com/hlandau/passlib"
	"golang.org/x/net/context"
)

const usersRootKey = "users"
const organizationsRootKey = "organizations"

// Store implements user data.Interface
type Store struct {
	registration string
	storage      storage.Interface
}

// NewStore returns an etcd implementation of user.Interface
func NewStore(storage storage.Interface, registration string) *Store {
	return &Store{storage: storage, registration: registration}
}

// Users

func (s *Store) rawUser(ctx context.Context, name string) (*User, error) {
	user := &User{}
	if err := s.storage.Get(ctx, path.Join(usersRootKey, name), user, true); err != nil {
		return nil, err
	}
	if user.GetName() == "" { // If there's no "name" in the answer, it means the user has not been found, so return nil
		return nil, nil
	}
	return user, nil
}

func secureUser(user *User) *User {
	if user == nil {
		return nil
	}
	user.PasswordHash = "" // For security reasons, remove the password hash
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

	// Check if email is already in use
	emailAlreadyUsed, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if emailAlreadyUsed != nil {
		return nil, EmailAlreadyUsed
	}

	// Check if organization with the same name already exists.
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
	if s.registration == configuration.RegistrationNone {
		user.IsVerified = true
	}
	if password, err = CheckPassword(password); err != nil {
		return nil, err
	}
	if user.PasswordHash, err = passlib.Hash(password); err != nil {
		return nil, err
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := s.storage.Create(ctx, path.Join(usersRootKey, name), user, nil, 0); err != nil {
		return nil, err
	}
	return secureUser(user), nil
}

// VerifyUser verifies a user account
func (s *Store) VerifyUser(ctx context.Context, userName string) error {
	// Update user
	uf := func(current proto.Message) (proto.Message, error) {
		user, ok := current.(*User)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected User): %T", user)
		}

		if user.TokenUsed {
			return nil, TokenAlreadyUsed
		}
		user.IsVerified = true
		user.TokenUsed = true
		return user, nil
	}
	if err := s.storage.Update(ctx, path.Join(usersRootKey, userName), uf, &User{}); err != nil {
		if err == storage.NotFound {
			return UserNotFound
		}
		return err
	}
	return nil
}

// CheckUserPassword checks the given user password
func (s *Store) CheckUserPassword(ctx context.Context, name string, password string) error {
	user, err := s.getVerifiedUser(ctx, name)
	if err != nil {
		return err
	}
	if _, err = passlib.Verify(password, user.PasswordHash); err != nil {
		return WrongPassword
	}
	return nil
}

// SetUserPassword sets the given user password
func (s *Store) SetUserPassword(ctx context.Context, name string, password string) error {
	// Password
	if _, err := CheckPassword(password); err != nil {
		return err
	}
	passwordHash, err := passlib.Hash(password)
	if err != nil {
		return err
	}

	// Update user
	uf := func(current proto.Message) (proto.Message, error) {
		user, ok := current.(*User)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected User): %T", user)
		}
		user.PasswordHash = passwordHash
		return user, nil
	}
	if err := s.storage.Update(ctx, path.Join(usersRootKey, name), uf, &User{}); err != nil {
		if err == storage.NotFound {
			return UserNotFound
		}
		return err
	}
	return nil
}

// GetUser fetches a user by name
func (s *Store) GetUser(ctx context.Context, name string) (user *User, err error) {
	if name, err = CheckName(name); err != nil {
		return nil, err
	}
	user, err = s.rawUser(ctx, name)
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

// GetUserOrganizations gets the organizations the given user is member of
func (s *Store) GetUserOrganizations(ctx context.Context, name string) ([]*Organization, error) {
	organizations, err := s.ListOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	userOrganizations := []*Organization{}
	for _, o := range organizations {
		if o.HasMember(name) {
			userOrganizations = append(userOrganizations, o)
		}
	}
	return userOrganizations, nil
}

// ListUsers lists users
func (s *Store) ListUsers(ctx context.Context) ([]*User, error) {
	protos := []proto.Message{}
	if err := s.storage.List(ctx, usersRootKey, storage.Everything, &User{}, &protos); err != nil {
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
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_USER, name}, DeleteAction, UserRN, name) {
		return NotAuthorized
	}

	// Get organizations this user is member of
	organizations, err := s.GetUserOrganizations(ctx, name)
	if err != nil {
		return err
	}
	// Check if user can be removed from all organizations
	for _, o := range organizations {
		if err := s.canRemoveUserFromOrganization(ctx, o.Name, name); err != nil {
			return err
		}
	}
	// If yes, remove the user from all organizations
	for _, o := range organizations {
		if err := s.RemoveUserFromOrganization(ctx, o.Name, name); err != nil {
			return err
		}
	}

	// Delete the user
	if err := s.storage.Delete(ctx, path.Join(usersRootKey, name), false, nil); err != nil {
		return err
	}
	return nil
}

// Organizations

func (s *Store) getOrganization(ctx context.Context, name string) (organization *Organization, err error) {
	if organization, err = s.GetOrganization(ctx, name); err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, OrganizationNotFound
	}
	return organization, nil
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
	if err := s.storage.Create(ctx, path.Join(organizationsRootKey, organization.Name), organization, nil, 0); err != nil {
		return err
	}
	return nil
}

// AddUserToOrganization adds a user to the given organization
func (s *Store) AddUserToOrganization(ctx context.Context, organizationName string, userName string) (err error) {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, OrganizationRN, organizationName) {
		return NotAuthorized
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Check if user is already a member
		if organization.HasMember(user.Name) {
			return nil, nil // User is already a member of the organization, return
		}

		// Add the user as a team member
		organization.Members = append(organization.Members, &OrganizationMember{
			Name: user.Name,
			Role: OrganizationRole_ORGANIZATION_MEMBER,
		})

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

func (s *Store) canRemoveUserFromOrganization(ctx context.Context, organizationName string, userName string) error {
	// Get organization
	organization, err := s.getOrganization(ctx, organizationName)
	if err != nil {
		return err
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Check if user is part of the organization
	memberIndex := organization.getMemberIndex(user.Name)
	if memberIndex == -1 {
		return nil // User is not a member of the organization, return
	}

	// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
	organization.Members = append(organization.Members[:memberIndex], organization.Members[memberIndex+1:]...)
	if err := organization.Validate(); err != nil {
		return err
	}
	return nil
}

// RemoveUserFromOrganization removes a user from the given organization
func (s *Store) RemoveUserFromOrganization(ctx context.Context, organizationName string, userName string) (err error) {
	// Check authorization
	if !(s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, OrganizationRN, organizationName) ||
		s.IsAuthorized(ctx, &Account{AccountType_USER, userName}, LeaveAction, OrganizationRN, organizationName)) {
		return NotAuthorized
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Check if user is part of the organization
		memberIndex := organization.getMemberIndex(user.Name)
		if memberIndex == -1 {
			return nil, nil // User is not a member of the organization, return
		}

		// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
		organization.Members = append(organization.Members[:memberIndex], organization.Members[memberIndex+1:]...)

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// ChangeOrganizationMemberRole changes the role of given user in the given organization
func (s *Store) ChangeOrganizationMemberRole(ctx context.Context, organizationName string, userName string, role OrganizationRole) (err error) {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, OrganizationRN, organizationName) {
		return NotAuthorized
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Check if user is already a member
		member := organization.getMember(userName)
		if member == nil {
			return nil, UserNotFound
		}

		// Change the role of the user
		member.Role = role

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// GetOrganization fetches a organization by name
func (s *Store) GetOrganization(ctx context.Context, name string) (organization *Organization, err error) {
	if name, err = CheckName(name); err != nil {
		return nil, err
	}
	organization = &Organization{}
	if err := s.storage.Get(ctx, path.Join(organizationsRootKey, name), organization, true); err != nil {
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
	if err := s.storage.List(ctx, organizationsRootKey, storage.Everything, &Organization{}, &protos); err != nil {
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
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, name}, DeleteAction, OrganizationRN, "") {
		return NotAuthorized
	}

	// Get organization
	_, err := s.getOrganization(ctx, name)
	if err != nil {
		return err
	}

	// Delete organization
	if err := s.storage.Delete(ctx, path.Join(organizationsRootKey, name), false, nil); err != nil {
		return err
	}
	return nil
}

// Teams

// CreateTeam creates a new team
func (s *Store) CreateTeam(ctx context.Context, organizationName, teamName string) error {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, CreateAction, TeamRN, teamName) {
		return NotAuthorized
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Check if team already exists
		if organization.hasTeam(teamName) {
			return nil, TeamAlreadyExists
		}

		// Create the new team
		team := &Team{
			Name:     teamName,
			CreateDt: time.Now().Unix(),
			Members:  []string{auth.GetUser(ctx)},
		}

		// Add the team to the organization
		organization.Teams = append(organization.Teams, team)

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// AddUserToTeam adds a user to the given team
func (s *Store) AddUserToTeam(ctx context.Context, organizationName string, teamName string, userName string) error {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, TeamRN, teamName) {
		return NotAuthorized
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Get team
		team := organization.getTeam(teamName)
		if team == nil {
			return nil, TeamNotFound
		}

		// Check if user is part of the organization
		if !organization.HasMember(user.Name) {
			return nil, NotPartOfOrganization
		}

		// Check if user is part of the team
		if team.hasMember(user.Name) {
			return nil, nil // User is already a member of the team, return
		}

		// Add the user as a team member
		team.Members = append(team.Members, user.Name)

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// RemoveUserFromTeam removes a user from the given team
func (s *Store) RemoveUserFromTeam(ctx context.Context, organizationName string, teamName string, userName string) error {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, TeamRN, teamName) {
		return NotAuthorized
	}

	// Get the user
	user, err := s.getVerifiedUser(ctx, userName)
	if err != nil {
		return err
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Get team
		team := organization.getTeam(teamName)
		if team == nil {
			return nil, TeamNotFound
		}

		// Check if user is actually a member
		memberIndex := team.getMemberIndex(user.Name)
		if memberIndex == -1 {
			return nil, nil // User is not a member of the team, return
		}

		// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
		team.Members = append(team.Members[:memberIndex], team.Members[memberIndex+1:]...)

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// AddResourceToTeam adds a resource to the given team
func (s *Store) AddResourceToTeam(ctx context.Context, organizationName string, teamName string, resourceID string) error {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, TeamRN, teamName) {
		return NotAuthorized
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Get team
		team := organization.getTeam(teamName)
		if team == nil {
			return nil, TeamNotFound
		}

		// Check if resource is already added to the team
		if team.hasResource(resourceID) {
			return nil, nil // Resource is already added to the team, return
		}

		// Add the resource
		team.Resources = append(team.Resources, &TeamResource{
			Id:              resourceID,
			PermissionLevel: TeamPermissionLevel_TEAM_READ,
		})

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// RemoveResourceFromTeam removes a resource from the given team
func (s *Store) RemoveResourceFromTeam(ctx context.Context, organizationName string, teamName string, resourceID string) error {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, TeamRN, teamName) {
		return NotAuthorized
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Get team
		team := organization.getTeam(teamName)
		if team == nil {
			return nil, TeamNotFound
		}

		// Check if resource is already present
		resourceIndex := team.getResourceIndex(resourceID)
		if resourceIndex == -1 {
			return nil, nil // Resource is not present, return
		}

		// Remove the resource from team resources. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
		team.Resources = append(team.Resources[:resourceIndex], team.Resources[resourceIndex+1:]...)

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// ChangeTeamResourcePermissionLevel changes the permission level over the given resource in the given team
func (s *Store) ChangeTeamResourcePermissionLevel(ctx context.Context, organizationName string, teamName string, resourceID string, permissionLevel TeamPermissionLevel) (err error) {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, OrganizationRN, organizationName) {
		return NotAuthorized
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Get team
		team := organization.getTeam(teamName)
		if team == nil {
			return nil, TeamNotFound
		}

		// Check if resource is associated to the team
		resource := team.getResourceById(resourceID)
		if resource == nil {
			return nil, ResourceNotFound
		}

		// Change the permission level over the given resource
		resource.PermissionLevel = permissionLevel

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// ChangeTeamName changes the name of given team
func (s *Store) ChangeTeamName(ctx context.Context, organizationName string, teamName, newName string) (err error) {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, UpdateAction, OrganizationRN, organizationName) {
		return NotAuthorized
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Check if team already exists
		alreadyExists := organization.getTeam(newName)
		if alreadyExists != nil {
			return nil, TeamAlreadyExists
		}

		// Get team
		team := organization.getTeam(teamName)
		if team == nil {
			return nil, TeamNotFound
		}

		// Update team name
		team.Name = newName

		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// GetTeam fetches a team by name
func (s *Store) GetTeam(ctx context.Context, organizationName string, teamName string) (*Team, error) {
	// Get organization
	organization, err := s.getOrganization(ctx, organizationName)
	if err != nil {
		return nil, err
	}

	// Get team
	if teamName, err = CheckName(teamName); err != nil {
		return nil, err
	}
	team := organization.getTeam(teamName)
	return team, nil
}

// ListTeams lists teams
func (s *Store) ListTeams(ctx context.Context, organizationName string) ([]*Team, error) {
	// Get organization
	organization, err := s.getOrganization(ctx, organizationName)
	if err != nil {
		return nil, err
	}
	return organization.Teams, nil
}

// DeleteTeam deletes a team by name
func (s *Store) DeleteTeam(ctx context.Context, organizationName string, teamName string) error {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_ORGANIZATION, organizationName}, DeleteAction, TeamRN, "") {
		return NotAuthorized
	}

	// Update organization
	uf := func(current proto.Message) (proto.Message, error) {
		organization, ok := current.(*Organization)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Organization): %T", organization)
		}

		// Check if the team is actually a team in the organization
		teamIndex := organization.getTeamIndex(teamName)
		if teamIndex == -1 {
			return nil, nil // Team is not part of the organization, return
		}

		// Remove the user from members. For details, check http://stackoverflow.com/questions/25025409/delete-element-in-a-slice
		organization.Teams = append(organization.Teams[:teamIndex], organization.Teams[teamIndex+1:]...)

		// Validate update
		if err := organization.Validate(); err != nil {
			return nil, err
		}
		return organization, nil
	}
	if err := s.storage.Update(ctx, path.Join(organizationsRootKey, organizationName), uf, &Organization{}); err != nil {
		if err == storage.NotFound {
			return OrganizationNotFound
		}
		return err
	}
	return nil
}

// Reset resets the account storage
func (s *Store) Reset(ctx context.Context) {
	s.storage.Delete(ctx, usersRootKey, true, nil)
	s.storage.Delete(ctx, organizationsRootKey, true, nil)
}
