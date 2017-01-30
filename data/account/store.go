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

// AccountUserByAltKey stores the alternate key value by name
const AccountUserByAltKey = AccountSchemaRootKey + "/account/name"

// AccountUserByIdKey key used to store the account protobuf type
const AccountUserByIdKey = AccountSchemaRootKey + "/account/id"

//AccountTeamByAltKey key used to store the alternate key by name
const AccountTeamByAltKey = AccountSchemaRootKey + "/team/name"

//AccountTeamByIdKey key used to store the team protobuf type
const AccountTeamByIdKey = AccountSchemaRootKey + "/team/id"

//AccountTeamByIdKey key used to store the team protobuf type
const AccountTeamMemberByIdKey = AccountSchemaRootKey + "/team/member/id"

//AccountResourceByIdKey key used to store the Resource protobuf type
const AccountResourceByIdKey = AccountSchemaRootKey + "/resource/id"

//AccountResourceByIdKey key used to store the Resource protobuf type
const AccountResourceByAltKey = AccountSchemaRootKey + "/resource/name"

//AccountResourceSettingsByIdKey key used to store the ResourceSettings protobuf type
const AccountResourceSettingsByIdKey = AccountSchemaRootKey + "/resource/settings/id"

//AccountPermissionByIdKey key used to store the ResourceSettings protobuf type
const AccountPermissionByIdKey = AccountSchemaRootKey + "/permission/id"

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

//AddResource adds a Resource to the resource table
func (s *Store) AddResource(resource *schema.Resource) (id string, err error) {

	// Integrity Check
	if err = s.checkResource(resource); err == nil {
		// Store the Account struct and the alternate key
		if err = s.Store.Create(s.ctx, path.Join(AccountResourceByIdKey, resource.Id), resource, nil, 0); err == nil {
			err = s.Store.Create(s.ctx, path.Join(AccountResourceByAltKey, resource.Name), &schema.Key{KeyId: resource.Id}, nil, 0)
		}
	}
	return resource.Id, err
}

func (s *Store) checkResource(resource *schema.Resource) (err error) {

	// Assign a new UUID if the resource does not exist
	if resource.Id == "" {
		resource.Id = generateUUID()

	} else {
		err = fmt.Errorf("Resource.Id must be blank/nil ")
	}
	return err
}

// GetResource returns a Resource from the Resource table
func (s *Store) GetResource(name string) (resource *schema.Resource, err error) {
	key := &schema.Key{}
	//Grab the ID
	err = s.Store.Get(s.ctx, path.Join(AccountResourceByAltKey, name), key, true)
	if err == nil && key.KeyId != "" {
		resource, err = s.GetResourceById(key.KeyId)
	}
	return resource, err
}

// GetResource returns a Resource from the Resource table
func (s *Store) GetResourceById(id string) (resource *schema.Resource, err error) {
	resource = &schema.Resource{}
	err = s.Store.Get(s.ctx, path.Join(AccountResourceByIdKey+"/"+id), resource, true)
	return resource, err
}

// GetResourceSettings returns a List of ResourceSettings from the ResourceSettings table
func (s *Store) GetResourceSettings(resourceId string) (rs []*schema.ResourceSettings, err error) {
	var out []proto.Message
	settings := &schema.ResourceSettings{}
	err = s.Store.List(s.ctx, AccountResourceSettingsByIdKey+"/"+resourceId, storage.Everything, settings, &out)
	if err == nil {
		// Unfortunately we have to iterate and filter
		for i := 0; i < len(out); i++ {
			m, ok := out[i].(*schema.ResourceSettings)
			if ok {
				rs = append(rs, m)
			} else if !ok {
				err = fmt.Errorf("Unexpected Type Encountered")
				break
			}
		}
	}
	return
}

//AddTeamMember adds a user account to the team table
func (s *Store) AddResourceSettings(rs *schema.ResourceSettings) (id string, err error) {

	if err = s.checkResourceSettings(rs); err == nil {
		// Store the Account struct and the alternate key
		err = s.Store.Create(s.ctx, path.Join(AccountResourceSettingsByIdKey, rs.ResourceId+"/"+rs.Id), rs, nil, 0)
	}
	return rs.Id, err
}

func (s *Store) checkResourceSettings(rs *schema.ResourceSettings) error {

	res, err := s.GetResource(rs.Id)
	if err == nil && rs.Id == "" {
		rs.Id = generateUUID()

	} else {
		err = fmt.Errorf("ResourceSetting %s already exists", res.Id)
	}
	return err
}

//AddTeamMember adds a user account to the team table
func (s *Store) AddTeamMember(member *schema.TeamMember) (id string, err error) {

	if err = s.checkTeamMember(member); err == nil {
		// Store the Account struct
		err = s.Store.Create(s.ctx, path.Join(AccountTeamMemberByIdKey, member.TeamId+"/"+member.Id), member, nil, 0)
	}
	return member.Id, err
}

func (s *Store) checkTeamMember(member *schema.TeamMember) error {

	mem, err := s.GetTeamMember(member.TeamId, member.Id)
	if err == nil && member.Id == "" {
		member.Id = generateUUID()

	} else {
		err = fmt.Errorf("TeamMember %s already exists", mem.Id)
	}
	return err
}

// GetTeamMember returns a TeamMember from the TeamMember table
func (s *Store) GetTeamMember(teamId string, memberId string) (member *schema.TeamMember, err error) {
	member = &schema.TeamMember{}
	key := &schema.Key{KeyId: memberId}
	err = s.Store.Get(s.ctx, path.Join(AccountTeamMemberByIdKey+"/"+teamId, key.KeyId), member, true)
	return member, err
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
			err = s.Store.Create(s.ctx, path.Join(AccountUserByAltKey, account.Name), &schema.Key{KeyId: account.Id}, nil, 0)
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
	if err = s.checkTeam(team); err == nil {
		// Store Team struct and alternate Key
		if err = s.Store.Create(s.ctx, path.Join(AccountTeamByIdKey, team.Id), team, nil, 0); err == nil {
			err = s.Store.Create(s.ctx, path.Join(AccountTeamByAltKey, team.Name), &schema.Key{KeyId: team.Id}, nil, 0)
		}
	}
	return team.Id, err
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

// GetTeam returns a Team from the Team table
func (s *Store) GetTeam(name string) (team *schema.Team, err error) {
	team = &schema.Team{}
	key := &schema.Key{}
	//Grab the ID
	err = s.Store.Get(s.ctx, path.Join(AccountTeamByAltKey, name), key, true)
	if err == nil && key.KeyId != "" {
		err = s.Store.Get(s.ctx, path.Join(AccountTeamByIdKey, key.KeyId), team, true)
	}
	return team, err
}

// GetAccount returns an account from the accounts table
func (s *Store) GetAccount(name string) (account *schema.Account, err error) {
	acct := &schema.Account{}
	key := &schema.Key{}
	//Grab the ID
	err = s.Store.Get(s.ctx, path.Join(AccountUserByAltKey, name), key, true)
	if err == nil && key.KeyId != "" {
		err = s.Store.Get(s.ctx, path.Join(AccountUserByIdKey, key.KeyId), acct, true)
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

//AddPermission adds a permission record to the Permission Table
func (s *Store) AddPermission(perm *schema.Permission) (id string, err error) {

	if err = s.checkPermission(perm); err == nil {
		// Store the Permission struct and the alternate key
		err = s.Store.Create(s.ctx, path.Join(AccountPermissionByIdKey, perm.ResourceId+"/"+perm.Id), perm, nil, 0)
	}
	return perm.Id, err
}

func (s *Store) checkPermission(perm *schema.Permission) (err error) {

	if perm.Id == "" {
		perm.Id = generateUUID()

	} else {
		err = fmt.Errorf("Perm ID must be blank/nil")
	}
	return err
}

//GetPermission retrieves a collection of permission records from the permission table
func (s *Store) GetPermission(resourceId string) (perms []*schema.Permission, err error) {
	var out []proto.Message
	perm := &schema.Permission{}
	err = s.Store.List(s.ctx, AccountPermissionByIdKey+"/"+perm.ResourceId, storage.Everything, perm, &out)
	if err == nil {
		// Unfortunately we have to iterate and filter
		for i := 0; i < len(out); i++ {
			m, ok := out[i].(*schema.Permission)
			if ok {
				perms = append(perms, m)
			} else if !ok {
				err = fmt.Errorf("Unexpected Type Encountered")
				break
			}
		}
	}
	return
}

//DeleteResourceSettings removes the Resource entry for a given Id
func (s *Store) DeleteResourceSettings(resourceId string) (err error) {

	err = s.Store.Delete(s.ctx, AccountResourceSettingsByIdKey+"/"+resourceId, true, nil)
	return err
}

//DeleteResource removes the Resource entry for a given Id
func (s *Store) DeleteResource(name string) (err error) {

	if resource, err := s.GetResource(name); err == nil && resource.Id != "" {
		//Cascade delete
		s.DeleteResourceSettings(resource.Id)
		err = s.Store.Delete(s.ctx, AccountResourceByIdKey+"/"+resource.Id, false, nil)
	}
	return
}

//DeleteTeamMember
func (s *Store) DeleteTeamMember(teamId string, memberId string) (err error) {
	//TODO
	return
}
