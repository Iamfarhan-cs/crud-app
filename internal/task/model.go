package task

import "time"

// Status describes where a task sits in its lifecycle.
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// Task is the domain model used inside the application.
// It can include server-owned fields and lifecycle fields because it represents
// the system's internal view of a task, not a client payload.
type Task struct {
	ID          string
	Title       string
	Description *string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
