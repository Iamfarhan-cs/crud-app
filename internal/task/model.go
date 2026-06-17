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

// CreateTaskRequest is separate from Task so clients can only send fields they
// are allowed to control when creating a task.
type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      *Status `json:"status,omitempty"`
}

// UpdateTaskRequest is separate from Task so partial updates can distinguish
// omitted fields from fields the client intentionally wants to change.
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *Status `json:"status,omitempty"`
}

// TaskResponse is separate from Task so API responses expose only the public
// representation and hide internal lifecycle fields such as DeletedAt.
type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
