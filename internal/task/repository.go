package task

import "context"

// Repository defines the persistence boundary for tasks.
// Implementations should store and retrieve task data.
// Future update and delete implementations must guard active rows with deleted_at IS NULL.
// HTTP concerns, business policy decisions, and database connection ownership must not live here.
type Repository interface {
	Create(ctx context.Context, task Task) (Task, error)

	// Active makes the soft-delete contract explicit: normal reads ignore tasks
	// whose deleted_at value is set.
	FindActiveByID(ctx context.Context, id string) (Task, error)

	// Active makes list behavior match the API contract: deleted tasks are hidden
	// from normal collection reads.
	ListActive(ctx context.Context) ([]Task, error)

	// Active makes update behavior explicit: deleted tasks cannot be modified.
	UpdateActive(ctx context.Context, task Task) (Task, error)

	SoftDelete(ctx context.Context, id string) error
}
