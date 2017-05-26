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
	// CreateStack creates a new stack
	CreateStack(ctx context.Context, name string) (stack *Stack, err error)

	// GetStack fetches a stack by id
	GetStack(ctx context.Context, id string) (stack *Stack, err error)

	// GetStackByName fetches a stack by name
	GetStackByName(ctx context.Context, name string) (stack *Stack, err error)

	// GetStackByFragmentOrName fetches a stack by fragment ID or name
	GetStackByFragmentOrName(ctx context.Context, fragmentOrName string) (stack *Stack, err error)

	// ListStacks lists stacks
	ListStacks(ctx context.Context) (stacks []*Stack, err error)

	// DeleteStack deletes a stack by id
	DeleteStack(ctx context.Context, id string) (err error)

	// Reset resets the stack storage
	Reset(ctx context.Context)
}
