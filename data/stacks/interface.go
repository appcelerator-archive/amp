package stacks

import "golang.org/x/net/context"

// Error type
type Error string

func (e Error) Error() string { return string(e) }

// Errors
const (
	InvalidName   = Error("name is invalid")
	AlreadyExists = Error("stack already exists")
	NotFound      = Error("stack not found")
)

// Interface defines the stack data access layer
type Interface interface {
	// Create creates a new stack
	Create(ctx context.Context, name string) (stack *Stack, err error)

	// Get fetches a stack by id
	Get(ctx context.Context, id string) (stack *Stack, err error)

	// GetByName fetches a stack by name
	GetByName(ctx context.Context, name string) (stack *Stack, err error)

	// GetByFragmentOrName fetches a stack by fragment ID or name
	GetByFragmentOrName(ctx context.Context, fragmentOrName string) (stack *Stack, err error)

	// List lists stacks
	List(ctx context.Context) (stacks []*Stack, err error)

	// Delete deletes a stack by id
	Delete(ctx context.Context, id string) (err error)
}
