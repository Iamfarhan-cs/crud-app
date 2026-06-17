package task

import "errors"

// Sentinel errors give upper layers stable failure meanings that can be mapped
// to HTTP responses by handlers.
// Wrapped errors must preserve errors.Is compatibility so callers can still
// detect these cases without depending on lower-level implementation details.
var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTitle      = errors.New("invalid title")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidPagination = errors.New("invalid pagination")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
)
