package dashboards

import (
	"fmt"
	"path"
	"time"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/go-docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const rootKey = "dashboards"

// Store implements dashboard data.Interface
type Store struct {
	storage  storage.Interface
	accounts accounts.Interface
}

// NewStore returns an etcd implementation of dashboards.Interface
func NewStore(storage storage.Interface, accounts accounts.Interface) *Store {
	return &Store{
		accounts: accounts,
		storage:  storage,
	}
}

// Create creates a new dashboard
func (s *Store) Create(ctx context.Context, name string, data string) (dashboard *Dashboard, err error) {
	// Check if dashboard already exists
	dashboardAlreadyExists, err := s.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if dashboardAlreadyExists != nil {
		return nil, AlreadyExists
	}

	// Create the new dashboard
	dashboard = &Dashboard{
		Id:       stringid.GenerateNonCryptoID(),
		Name:     name,
		Owner:    accounts.GetRequesterAccount(ctx),
		CreateDt: time.Now().Unix(),
		Data:     data,
	}
	if err := dashboard.Validate(); err != nil {
		return nil, err
	}
	if err := s.storage.Create(ctx, path.Join(rootKey, dashboard.Id), dashboard, nil, 0); err != nil {
		return nil, err
	}
	return dashboard, nil
}

// Get fetches a dashboard by id
func (s *Store) Get(ctx context.Context, id string) (*Dashboard, error) {
	dashboard := &Dashboard{}
	if err := s.storage.Get(ctx, path.Join(rootKey, id), dashboard, true); err != nil {
		return nil, err
	}
	// If there's no "id" in the answer, it means the dashboard has not been found, so return nil
	if dashboard.GetId() == "" {
		return nil, nil
	}
	return dashboard, nil
}

// GetByName fetches a dashboard by name
func (s *Store) GetByName(ctx context.Context, name string) (dashboard *Dashboard, err error) {
	if name, err = CheckName(name); err != nil {
		return nil, err
	}
	dashboards, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, dashboard := range dashboards {
		if dashboard.Name == name {
			return dashboard, nil
		}
	}
	return nil, nil
}

// List lists dashboards
func (s *Store) List(ctx context.Context) ([]*Dashboard, error) {
	protos := []proto.Message{}
	if err := s.storage.List(ctx, rootKey, storage.Everything, &Dashboard{}, &protos); err != nil {
		return nil, err
	}
	dashboards := []*Dashboard{}
	for _, proto := range protos {
		dashboard := proto.(*Dashboard)
		if !s.accounts.IsAuthorized(ctx, dashboard.Owner, accounts.ReadAction, accounts.DashboardRN, dashboard.Id) {
			continue
		}

		// Check if we have an active organization
		switch accounts.GetRequesterAccount(ctx).Organization {
		case "": // If there's no active organization, add the dashboard to the results
			dashboards = append(dashboards, dashboard)
		case dashboard.Owner.Organization: // If the dashboard belongs to the active organization, add the dashboard to the results
			dashboards = append(dashboards, dashboard)
		case accounts.SuperOrganization: // If the requester is a member of the super organization, add the stack to the results
			dashboards = append(dashboards, dashboard)
		default:
		}
	}
	return dashboards, nil
}

// UpdateName renames the given dashboard
func (s *Store) UpdateName(ctx context.Context, id string, name string) error {
	uf := func(current proto.Message) (proto.Message, error) {
		dashboard, ok := current.(*Dashboard)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Dashboard): %T", dashboard)
		}
		dashboard.Name = name
		return dashboard, nil
	}
	if err := s.storage.Update(ctx, path.Join(rootKey, id), uf, &Dashboard{}); err != nil {
		if err == storage.NotFound {
			return NotFound
		}
		return err
	}
	return nil
}

// UpdateData updates the given dashboard data
func (s *Store) UpdateData(ctx context.Context, id string, data string) error {
	uf := func(current proto.Message) (proto.Message, error) {
		dashboard, ok := current.(*Dashboard)
		if !ok {
			return nil, fmt.Errorf("value is not the right type (expected Dashboard): %T", dashboard)
		}
		dashboard.Data = data
		return dashboard, nil
	}
	if err := s.storage.Update(ctx, path.Join(rootKey, id), uf, &Dashboard{}); err != nil {
		if err == storage.NotFound {
			return NotFound
		}
		return err
	}
	return nil
}

// Delete deletes a dashboard by id
func (s *Store) Delete(ctx context.Context, id string) error {
	dashboard, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if dashboard == nil {
		return NotFound
	}

	// Check authorization
	if !s.accounts.IsAuthorized(ctx, dashboard.Owner, accounts.DeleteAction, accounts.DashboardRN, dashboard.Id) {
		return accounts.NotAuthorized
	}

	// Delete the dashboard in all teams of the owning organization
	if dashboard.Owner.Organization != "" {
		org, err := s.accounts.GetOrganization(ctx, dashboard.Owner.Organization)
		if err != nil {
			return err
		}
		if org == nil {
			return accounts.OrganizationNotFound
		}
		for _, team := range org.Teams {
			s.accounts.RemoveResourceFromTeam(ctx, dashboard.Owner.Organization, team.Name, dashboard.Id)
		}
	}

	// Delete the dashboard
	if err := s.storage.Delete(ctx, path.Join(rootKey, dashboard.Id), false, nil); err != nil {
		return err
	}
	return nil
}

// Reset resets the account storage
func (s *Store) Reset(ctx context.Context) {
	s.storage.Delete(ctx, rootKey, true, nil)
}
