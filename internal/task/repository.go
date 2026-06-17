package task

import "context"

// Repository is the persistence boundary for tasks.
// The service depends on this contract, not on PostgreSQL or any other storage
// implementation. Concrete repositories are responsible for data access only.
// Repository implementations must use parameterized SQL only and must never
// concatenate user input into query strings.
// Repository code must not contain HTTP concerns such as status codes, request
// parsing, response formatting, or route behavior.
// Business rules and business-policy decisions belong in the service layer.
type Repository interface {
	Create(ctx context.Context, task Task) (Task, error)

	// Active means deleted_at IS NULL. Normal reads exclude soft-deleted tasks.
	FindActiveByID(ctx context.Context, id string) (Task, error)

	// Active means deleted_at IS NULL. Normal list reads exclude soft-deleted
	// tasks, and limit/offset keep list queries bounded.
	ListActive(ctx context.Context, limit int, offset int) ([]Task, error)

	// Active means deleted_at IS NULL. Updates must not modify soft-deleted tasks.
	UpdateActive(ctx context.Context, task Task) (Task, error)

	// SoftDelete performs a soft delete by marking the task deleted instead of
	// physically removing it.
	SoftDelete(ctx context.Context, id string) error
}
