package task

import "errors"

// Domain errors give upper layers stable failure meanings.
// This file should contain task-specific error values.
// HTTP status codes, SQL driver errors, and response bodies must not live here.
var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidTask  = errors.New("invalid task")
)
