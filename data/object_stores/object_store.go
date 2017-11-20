package object_stores

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const rootKey = "stores"

// Store implements ObjectStore data.Interface
type Store struct {
	accounts accounts.Interface
	storage  storage.Interface
}

// NewStore returns an etcd implementation of ObjectStore.Interface
func NewStore(storage storage.Interface, accounts accounts.Interface) *Store {
	return &Store{
		accounts: accounts,
		storage:  storage,
	}
}

// ObjectStores

func (s *Store) isNameAvailable(ctx context.Context, name string) (bool, error) {
	protos := []proto.Message{}
	if err := s.storage.List(ctx, rootKey, storage.Everything, &ObjectStore{}, &protos); err != nil {
		return false, err
	}
	for _, proto := range protos {
		if strings.EqualFold(proto.(*ObjectStore).Name, name) {
			return false, nil
		}
	}
	return true, nil
}

// Create creates a new objectStore
func (s *Store) Create(ctx context.Context, name string) (ostore *ObjectStore, err error) {
	// Check if ostore already exists
	nameAvailable, err := s.isNameAvailable(ctx, name)
	if err != nil {
		return nil, err
	}
	if !nameAvailable {
		return nil, AlreadyExists
	}

	// Create the new ostore
	ostore = &ObjectStore{
		Id:       stringid.GenerateNonCryptoID(),
		Name:     name,
		Owner:    accounts.GetRequesterAccount(ctx),
		CreateDt: time.Now().Unix(),
	}
	if err := ostore.Validate(); err != nil {
		return nil, err
	}
	if err := s.storage.Create(ctx, path.Join(rootKey, ostore.Id), ostore, nil, 0); err != nil {
		return nil, err
	}
	return ostore, nil
}

// Get fetches an object store by id
func (s *Store) Get(ctx context.Context, id string) (*ObjectStore, error) {
	ostore := &ObjectStore{}
	if err := s.storage.Get(ctx, path.Join(rootKey, id), ostore, true); err != nil {
		return nil, err
	}
	// If there's no "id" in the answer, it means the object store has not been found, so return nil
	if ostore.GetId() == "" {
		return nil, nil
	}
	return ostore, nil
}

// GetByName fetches an object store by name
func (s *Store) GetByName(ctx context.Context, name string) (ostore *ObjectStore, err error) {
	if name, err = CheckName(name); err != nil {
		return nil, err
	}
	ostores, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, ostore := range ostores {
		if ostore.Name == name {
			return ostore, nil
		}
	}
	return nil, nil
}

// GetByFragmentOrName fetches an object store by fragment ID or name
func (s *Store) GetByFragmentOrName(ctx context.Context, fragmentOrName string) (ostore *ObjectStore, err error) {
	objs, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if obj.Name == fragmentOrName || strings.HasPrefix(strings.ToLower(obj.Id), strings.ToLower(fragmentOrName)) {
			ostore = obj
			break
		}
	}
	return ostore, nil
}

// List lists object stores
func (s *Store) List(ctx context.Context) ([]*ObjectStore, error) {
	protos := []proto.Message{}
	if err := s.storage.List(ctx, rootKey, storage.Everything, &ObjectStore{}, &protos); err != nil {
		return nil, err
	}
	ostores := []*ObjectStore{}
	for _, proto := range protos {
		ostore := proto.(*ObjectStore)
		if !s.accounts.IsAuthorized(ctx, ostore.Owner, accounts.ReadAction, accounts.ObjectStoreRN, ostore.Id) {
			continue
		}

		// Check if we have an active organization
		switch accounts.GetRequesterAccount(ctx).Organization {
		case "": // If there's no active organization, add the object to the results
			ostores = append(ostores, ostore)
		case ostore.Owner.Organization: // If the object store belongs to the active organization, add the object store to the results
			ostores = append(ostores, ostore)
		case accounts.SuperOrganization: // If the requester is a member of the super organization, add the object store to the results
			ostores = append(ostores, ostore)
		default:
		}
	}
	return ostores, nil
}

// Delete deletes an object store by id
func (s *Store) Delete(ctx context.Context, id string) error {
	ostore, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if ostore == nil {
		return NotFound
	}

	// Check authorization
	if !s.accounts.IsAuthorized(ctx, ostore.Owner, accounts.DeleteAction, accounts.ObjectStoreRN, ostore.Id) {
		return accounts.NotAuthorized
	}

	// Delete the object store in all teams of the owning organization
	if ostore.Owner.Organization != "" {
		org, err := s.accounts.GetOrganization(ctx, ostore.Owner.Organization)
		if err != nil {
			return err
		}
		if org == nil {
			return accounts.OrganizationNotFound
		}
		for _, team := range org.Teams {
			s.accounts.RemoveResourceFromTeam(ctx, ostore.Owner.Organization, team.Name, ostore.Id)
		}
	}

	// Delete the object store
	if err := s.storage.Delete(ctx, path.Join(rootKey, ostore.Id), false, nil); err != nil {
		return err
	}
	return nil
}

// UpdateLocation updates the location of the object store
func (s *Store) UpdateLocation(ctx context.Context, id string, location string) error {
	uf := func(current proto.Message) (proto.Message, error) {
		ostore, ok := current.(*ObjectStore)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected ObjectStore): %T", ostore)
		}
		ostore.Location = location
		return ostore, nil
	}
	if err := s.storage.Update(ctx, path.Join(rootKey, id), uf, &ObjectStore{}); err != nil {
		if err == storage.NotFound {
			return NotFound
		}
		return err
	}
	return nil
}
