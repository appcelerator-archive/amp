package dashboards

import "golang.org/x/net/context"

// Error type
type Error string

func (e Error) Error() string { return string(e) }

// Errors
const (
	InvalidName   = Error("name is invalid")
	AlreadyExists = Error("dashboard already exists")
	NotFound      = Error("dashboard not found")
)

// Interface defines the dashboard data access layer
type Interface interface {
	// Create creates a new dashboard
	Create(ctx context.Context, name string, data string) (dashboard *Dashboard, err error)

	// Get fetches a dashboard by id
	Get(ctx context.Context, id string) (dashboard *Dashboard, err error)

	// GetByName fetches a dashboard by name
	GetByName(ctx context.Context, name string) (dashboard *Dashboard, err error)

	// List lists dashboards
	List(ctx context.Context) (dashboards []*Dashboard, err error)

	// UpdateName renames the given dashboard
	UpdateName(ctx context.Context, id string, name string) (err error)

	// UpdateData updates the given dashboard data
	UpdateData(ctx context.Context, id string, data string) (err error)

	// Delete deletes a dashboard by id
	Delete(ctx context.Context, id string) (err error)

	// Reset resets the dashboard storage
	Reset(ctx context.Context)
}
