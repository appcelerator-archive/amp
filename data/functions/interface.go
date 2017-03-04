package functions

import (
	"context"
)

type Error string

func (e Error) Error() string { return string(e) }

const InvalidName = Error("name is invalid")
const InvalidImage = Error("image is invalid")
const FunctionAlreadyExists = Error("function already exists")

// Interface defines the function data access layer
type Interface interface {
	// CreateFunction creates a new function
	CreateFunction(ctx context.Context, name string, image string) (function *Function, err error)

	// GetFunction fetches a function by id
	GetFunction(ctx context.Context, id string) (function *Function, err error)

	// GetFunctionByName fetches a function by name
	GetFunctionByName(ctx context.Context, name string) (function *Function, err error)

	// ListFunctions lists functions
	ListFunctions(ctx context.Context) (functions []*Function, err error)

	// DeleteFunction deletes a function by id
	DeleteFunction(ctx context.Context, id string) (err error)

	// Reset resets the function store
	Reset(ctx context.Context)
}
