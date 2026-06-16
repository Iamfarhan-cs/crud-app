package task

// Service owns task business rules and use-case orchestration.
// It should validate domain behavior and coordinate repository calls.
// Future concurrency rules such as stale update conflicts and idempotent create behavior belong here.
// HTTP request parsing, response formatting, SQL, and database connection setup must not live here.
type Service struct {
	repository Repository
}

// NewService prepares the task service with its persistence dependency.
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
