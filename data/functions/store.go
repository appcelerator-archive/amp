package functions

import (
	"path"
	"time"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const functionsRootKey = "functions"

// Store implements function data.Interface
type Store struct {
	store    storage.Interface
	accounts accounts.Interface
}

// NewStore returns an etcd implementation of function.Interface
func NewStore(store storage.Interface) *Store {
	return &Store{
		store:    store,
		accounts: accounts.NewStore(store),
	}
}

// Functions

// CreateFunction creates a new function
func (s *Store) CreateFunction(ctx context.Context, name string, image string) (function *Function, err error) {
	// Check if function already exists
	functionAlreadyExists, err := s.GetFunctionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if functionAlreadyExists != nil {
		return nil, FunctionAlreadyExists
	}

	// Create the new function
	function = &Function{
		Id:       stringid.GenerateNonCryptoID(),
		Name:     name,
		Image:    image,
		Owner:    accounts.GetRequesterAccount(ctx),
		CreateDt: time.Now().Unix(),
	}
	if err := function.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.Create(ctx, path.Join(functionsRootKey, function.Id), function, nil, 0); err != nil {
		return nil, err
	}
	return function, nil
}

// GetFunction fetches a function by id
func (s *Store) GetFunction(ctx context.Context, id string) (*Function, error) {
	function := &Function{}
	if err := s.store.Get(ctx, path.Join(functionsRootKey, id), function, true); err != nil {
		return nil, err
	}
	// If there's no "id" in the answer, it means the function has not been found, so return nil
	if function.GetId() == "" {
		return nil, nil
	}
	return function, nil
}

// GetFunctionByName fetches a function by name
func (s *Store) GetFunctionByName(ctx context.Context, name string) (*Function, error) {
	if err := CheckName(name); err != nil {
		return nil, err
	}
	functions, err := s.ListFunctions(ctx)
	if err != nil {
		return nil, err
	}
	for _, function := range functions {
		if function.Name == name {
			return function, nil
		}
	}
	return nil, nil
}

// ListFunctions lists functions
func (s *Store) ListFunctions(ctx context.Context) ([]*Function, error) {
	protos := []proto.Message{}
	if err := s.store.List(ctx, functionsRootKey, storage.Everything, &Function{}, &protos); err != nil {
		return nil, err
	}
	functions := []*Function{}
	for _, proto := range protos {
		functions = append(functions, proto.(*Function))
	}
	return functions, nil
}

// DeleteFunction deletes a function by id
func (s *Store) DeleteFunction(ctx context.Context, id string) error {
	f, err := s.GetFunction(ctx, id)
	if err != nil {
		return err
	}
	if f == nil {
		return FunctionNotFound
	}

	// Check authorization
	if !s.accounts.IsAuthorized(ctx, f.Owner, accounts.DeleteAction, accounts.FunctionResource) {
		return accounts.NotAuthorized
	}

	// Delete the function
	if err := s.store.Delete(ctx, path.Join(functionsRootKey, id), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account store
func (s *Store) Reset(ctx context.Context) {
	s.store.Delete(ctx, functionsRootKey, true, nil)
}
