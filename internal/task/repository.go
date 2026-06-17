package task

import "context"

// Repository is the persistence boundary for tasks.
// The service depends on this contract, not on PostgreSQL or any other storage
// implementation. Concrete repositories are responsible for data access only.
// Repository implementations must use parameterized SQL and must never
// concatenate user input into query strings.
// Repository code must not contain HTTP concerns such as status codes, request
// parsing, response formatting, or route behavior.
// Repository code must not make business-policy decisions; those belong in the
// service layer.
type Repository interface {
	Create(ctx context.Context, task Task) (Task, error)

	// Active means deleted_at IS NULL, so normal reads ignore soft-deleted tasks.
	FindActiveByID(ctx context.Context, id string) (Task, error)

	// Active means deleted_at IS NULL, and limit/offset keep list queries bounded.
	ListActive(ctx context.Context, limit int, offset int) ([]Task, error)

	// Active means deleted_at IS NULL, so deleted tasks cannot be updated.
	UpdateActive(ctx context.Context, task Task) (Task, error)

	// SoftDelete marks a task as deleted instead of physically removing it.
	SoftDelete(ctx context.Context, id string) error
}
