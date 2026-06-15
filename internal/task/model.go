package task

import "time"

// Task is the core domain entity for the task management API.
// This file should contain task data shapes and domain-level names.
// It must not contain HTTP parsing, database queries, or service orchestration.
type Task struct {
	ID          string
	Title       string
	Description string
	Status      Status
	Priority    Priority
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Status describes where a task sits in its lifecycle.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// Priority describes the relative importance of a task.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)
