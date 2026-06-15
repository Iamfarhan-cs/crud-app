package task

import "time"

// CreateTaskRequest is the client payload for creating a task.
type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Status      *Status   `json:"status,omitempty"`
	Priority    *Priority `json:"priority,omitempty"`
	DueDate     *string   `json:"due_date,omitempty"`
}

// UpdateTaskRequest is the client payload for partially updating a task.
type UpdateTaskRequest struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Status      *Status   `json:"status,omitempty"`
	Priority    *Priority `json:"priority,omitempty"`
	DueDate     *string   `json:"due_date,omitempty"`
}

// TaskResponse is the API representation of a task.
type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Status      Status    `json:"status"`
	Priority    Priority  `json:"priority"`
	DueDate     *string   `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListTasksResponse is the paginated response for task collections.
type ListTasksResponse struct {
	Data       []TaskResponse `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

// ErrorResponse is the standard error body for failed API requests.
type ErrorResponse struct {
	Error struct {
		Code    string         `json:"code"`
		Message string         `json:"message"`
		Details map[string]any `json:"details,omitempty"`
	} `json:"error"`
}
