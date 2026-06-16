package task

import "net/http"

// Handler owns HTTP concerns for task routes.
// It should translate requests into service calls and service results into HTTP responses.
// Handler-level security work includes request parsing, request size assumptions,
// safe client error responses, and basic input shape validation.
// Authentication, authorization, business rules, persistence details, and database connection setup must not live here.
type Handler struct {
	service *Service
}

// NewHandler prepares task HTTP handlers for future route registration.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes is the future home for task route registration.
// Complete endpoint behavior will be added after the architecture boundaries are stable.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	_ = h
	_ = mux
}
