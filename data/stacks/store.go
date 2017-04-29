package stacks

import (
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const stacksRootKey = "stacks"

// Store implements stack data.Interface
type Store struct {
	store    storage.Interface
	accounts accounts.Interface
}

// NewStore returns an etcd implementation of stacks.Interface
func NewStore(store storage.Interface) *Store {
	return &Store{
		store:    store,
		accounts: accounts.NewStore(store),
	}
}

// Stacks

// CreateStack creates a new stack
func (s *Store) CreateStack(ctx context.Context, name string) (stack *Stack, err error) {
	// Check if stack already exists
	stackAlreadyExists, err := s.GetStackByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if stackAlreadyExists != nil {
		return nil, StackAlreadyExists
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
	if err := s.store.Create(ctx, path.Join(stacksRootKey, stack.Id), stack, nil, 0); err != nil {
		return nil, err
	}
	return stack, nil
}

// GetStack fetches a stack by id
func (s *Store) GetStack(ctx context.Context, id string) (*Stack, error) {
	stack := &Stack{}
	if err := s.store.Get(ctx, path.Join(stacksRootKey, id), stack, true); err != nil {
		return nil, err
	}
	// If there's no "id" in the answer, it means the stack has not been found, so return nil
	if stack.GetId() == "" {
		return nil, nil
	}
	return stack, nil
}

// GetStackByName fetches a stack by name
func (s *Store) GetStackByName(ctx context.Context, name string) (stack *Stack, err error) {
	if name, err = CheckName(name); err != nil {
		return nil, err
	}
	stacks, err := s.ListStacks(ctx)
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

// GetStackByFragmentOrName fetches a stack by fragment ID or name
func (s *Store) GetStackByFragmentOrName(ctx context.Context, fragmentOrName string) (stack *Stack, err error) {
	stks, err := s.ListStacks(ctx)
	if err != nil {
		return nil, err
	}
	for _, stk := range stks {
		if stk.Name == fragmentOrName || strings.HasPrefix(strings.ToLower(stk.Id), strings.ToLower(fragmentOrName)) {
			stack = stk
			break
		}
	}
	if stack == nil {
		return nil, StackNotFound
	}
	return stack, nil
}

// ListStacks lists stacks
func (s *Store) ListStacks(ctx context.Context) ([]*Stack, error) {
	protos := []proto.Message{}
	if err := s.store.List(ctx, stacksRootKey, storage.Everything, &Stack{}, &protos); err != nil {
		return nil, err
	}
	stacks := []*Stack{}
	for _, proto := range protos {
		stacks = append(stacks, proto.(*Stack))
	}
	return stacks, nil
}

// DeleteStack deletes a stack by id
func (s *Store) DeleteStack(ctx context.Context, id string) error {
	stack, err := s.GetStack(ctx, id)
	if err != nil {
		return err
	}
	if stack == nil {
		return StackNotFound
	}

	// Check authorization
	if !s.accounts.IsAuthorized(ctx, stack.Owner, accounts.DeleteAction, accounts.StackRN, stack.Id) {
		return accounts.NotAuthorized
	}

	// Delete the stack
	if err := s.store.Delete(ctx, path.Join(stacksRootKey, stack.Id), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) {
	s.store.Delete(ctx, stacksRootKey, true, nil)
}
