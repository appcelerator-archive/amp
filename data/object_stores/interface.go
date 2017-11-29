package object_stores

import "golang.org/x/net/context"

// Error type
type Error string

func (e Error) Error() string { return string(e) }

// Errors
const (
	InvalidName       = Error("name is invalid")
	AlreadyExists     = Error("object store already exists, in this account or another")
	NotFound          = Error("object store not found")
	NotImplemented    = Error("object store not implemented for this provider")
	AlreadyOwnedByYou = Error("object store already exists on this account")
)

// Interface defines the object store data access layer
type Interface interface {
	// Create creates a new object store
	Create(ctx context.Context, name string) (store *ObjectStore, err error)

	// Get fetches an object store by id
	Get(ctx context.Context, id string) (store *ObjectStore, err error)

	// GetByName fetches an object store by name
	GetByName(ctx context.Context, name string) (store *ObjectStore, err error)

	// GetByFragmentOrName fetches an object store by fragment ID or name
	GetByFragmentOrName(ctx context.Context, fragmentOrName string) (store *ObjectStore, err error)

	// List lists object stores
	List(ctx context.Context) (stores []*ObjectStore, err error)

	// Delete deletes an object store by id
	Delete(ctx context.Context, id string) (err error)

	// UpdateLocation updates the location of the object store
	UpdateLocation(ctx context.Context, id string, location string) (err error)
}
