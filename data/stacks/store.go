package stacks

import (
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const rootKey = "stacks"

// Store implements stack data.Interface
type Store struct {
	accounts accounts.Interface
	storage  storage.Interface
}

// NewStore returns an etcd implementation of stacks.Interface
func NewStore(storage storage.Interface, accounts accounts.Interface) *Store {
	return &Store{
		accounts: accounts,
		storage:  storage,
	}
}

// Stacks

func (s *Store) isNameAvailable(ctx context.Context, name string) (bool, error) {
	protos := []proto.Message{}
	if err := s.storage.List(ctx, rootKey, storage.Everything, &Stack{}, &protos); err != nil {
		return false, err
	}
	for _, proto := range protos {
		if strings.EqualFold(proto.(*Stack).Name, name) {
			return false, nil
		}
	}
	return true, nil
}

// Create creates a new stack
func (s *Store) Create(ctx context.Context, name string) (stack *Stack, err error) {
	// Check if stack already exists
	nameAvailable, err := s.isNameAvailable(ctx, name)
	if err != nil {
		return nil, err
	}
	if !nameAvailable {
		return nil, AlreadyExists
	}

	// Create the new stack
	stack = &Stack{
		Id:       stringid.GenerateNonCryptoID(),
		Name:     name,
		Owner:    accounts.GetRequesterAccount(ctx),
		CreateDt: time.Now().Unix(),
	}
	if err := stack.Validate(); err != nil {
		return nil, err
	}
	if err := s.storage.Create(ctx, path.Join(rootKey, stack.Id), stack, nil, 0); err != nil {
		return nil, err
	}
	return stack, nil
}

// Get fetches a stack by id
func (s *Store) Get(ctx context.Context, id string) (*Stack, error) {
	stack := &Stack{}
	if err := s.storage.Get(ctx, path.Join(rootKey, id), stack, true); err != nil {
		return nil, err
	}
	// If there's no "id" in the answer, it means the stack has not been found, so return nil
	if stack.GetId() == "" {
		return nil, nil
	}
	return stack, nil
}

// GetByName fetches a stack by name
func (s *Store) GetByName(ctx context.Context, name string) (stack *Stack, err error) {
	if name, err = CheckName(name); err != nil {
		return nil, err
	}
	stacks, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, stack := range stacks {
		if stack.Name == name {
			return stack, nil
		}
	}
	return nil, nil
}

// GetByFragmentOrName fetches a stack by fragment ID or name
func (s *Store) GetByFragmentOrName(ctx context.Context, fragmentOrName string) (stack *Stack, err error) {
	stks, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, stk := range stks {
		if stk.Name == fragmentOrName || strings.HasPrefix(strings.ToLower(stk.Id), strings.ToLower(fragmentOrName)) {
			stack = stk
			break
		}
	}
	return stack, nil
}

// List lists stacks
func (s *Store) List(ctx context.Context) ([]*Stack, error) {
	protos := []proto.Message{}
	if err := s.storage.List(ctx, rootKey, storage.Everything, &Stack{}, &protos); err != nil {
		return nil, err
	}
	stacks := []*Stack{}
	for _, proto := range protos {
		stack := proto.(*Stack)
		if !s.accounts.IsAuthorized(ctx, stack.Owner, accounts.ReadAction, accounts.StackRN, stack.Id) {
			continue
		}

		// Check if we have an active organization
		switch accounts.GetRequesterAccount(ctx).Organization {
		case "": // If there's no active organization, add the stack to the results
			stacks = append(stacks, stack)
		case stack.Owner.Organization: // If the stack belongs to the active organization, add the stack to the results
			stacks = append(stacks, stack)
		case accounts.SuperOrganization: // If the requester is a member of the super organization, add the stack to the results
			stacks = append(stacks, stack)
		default:
		}
	}
	return stacks, nil
}

// Delete deletes a stack by id
func (s *Store) Delete(ctx context.Context, id string) error {
	stack, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if stack == nil {
		return NotFound
	}

	// Check authorization
	if !s.accounts.IsAuthorized(ctx, stack.Owner, accounts.DeleteAction, accounts.StackRN, stack.Id) {
		return accounts.NotAuthorized
	}

	// Delete the stack in all teams of the owning organization
	if stack.Owner.Organization != "" {
		org, err := s.accounts.GetOrganization(ctx, stack.Owner.Organization)
		if err != nil {
			return err
		}
		if org == nil {
			return accounts.OrganizationNotFound
		}
		for _, team := range org.Teams {
			s.accounts.RemoveResourceFromTeam(ctx, stack.Owner.Organization, team.Name, stack.Id)
		}
	}

	// Delete the stack
	if err := s.storage.Delete(ctx, path.Join(rootKey, stack.Id), false, nil); err != nil {
		return err
	}
	return nil
}
