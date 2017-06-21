package accounts

import (
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/hlandau/passlib"
	"github.com/ory/ladon"
	"github.com/ory/ladon/manager/memory"
	"golang.org/x/net/context"
)

const superAccountRootKey = "sa"
const usersRootKey = "users"
const organizationsRootKey = "organizations"
const superUser = "su"
const superOrganization = "so"

// Store implements user data.Interface
type Store struct {
	registration string
	storage      storage.Interface
	warden       *ladon.Ladon
}

// NewStore returns a new accounts storage
func NewStore(s storage.Interface, registration string, SUPassword string) (*Store, error) {
	store := &Store{
		storage:      s,
		registration: registration,
		warden: &ladon.Ladon{
			Manager: memory.NewMemoryManager(),
		},
	}

	// Register policies
	for _, policy := range policies {
		if err := store.warden.Manager.Create(policy); err != nil {
			log.Fatal("Unable to create policy:", err)
		}
	}

	if err := store.createSuperAccounts(SUPassword); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) createSuperAccounts(SUPassword string) error {
	if SUPassword == "" {
		log.Println("SUPassword is empty. Skipping creation of super accounts.")
		return nil
	}

	// Add a policy giving full access to super organization members
	s.warden.Manager.Create(&ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{"<.*>"},
		Actions:   []string{"<.*>"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &SuperUserCondition{},
		},
	})

	// Check if super accounts haven been created already
	ctx := context.Background()
	err := s.storage.Create(ctx, path.Join(superAccountRootKey, "created"), &User{}, nil, 0)
	switch err {
	case nil: // No error, do nothing
	case storage.AlreadyExists: // Super accounts have already been created, just return
		log.Println("Super accounts already created")
		return nil
	default:
		return err
	}

	// Create the initial super user
	user, err := s.GetUser(ctx, superUser)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("Initial super user should not exist already. Check the storage.")
	}
	su := &User{
		Name:       superUser,
		Email:      "super@user.amp",
		IsVerified: true,
		CreateDt:   time.Now().Unix(),
	}
	if su.PasswordHash, err = passlib.Hash(SUPassword); err != nil {
		return err
	}
	if err := s.storage.Create(ctx, path.Join(usersRootKey, su.Name), su, nil, 0); err != nil {
		return err
	}
	log.Println("Successfully created initial super user")
	return nil
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

// DeleteNotVerifiedUser deletes the user by name only if it's not verified
func (s *Store) DeleteNotVerifiedUser(ctx context.Context, name string) error {

	//Get user to verify it is well not verified
	user, err := s.GetUser(ctx, name)
	if err != nil {
		return err
	}
	if user != nil && !user.IsVerified {
		if err := s.storage.Delete(ctx, path.Join(usersRootKey, name), false, nil); err != nil {
			return err
		}
	}
	return nil
}

// DeleteUser deletes a user by name
func (s *Store) DeleteUser(ctx context.Context, name string) error {
	// Check authorization
	if !s.IsAuthorized(ctx, &Account{AccountType_USER, name}, DeleteAction, UserRN, name) {
		return NotAuthorized
	}

	// Delete the user
	if err := s.storage.Delete(ctx, path.Join(usersRootKey, name), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account storage
func (s *Store) Reset(ctx context.Context) {
	s.storage.Delete(ctx, usersRootKey, true, nil)
	s.storage.Delete(ctx, organizationsRootKey, true, nil)
}
