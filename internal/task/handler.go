package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const maxRequestBodySize = 1 << 20

// Handler owns HTTP concerns for task routes.
// It translates HTTP requests into service calls and service results into HTTP
// responses. Business rules, SQL, and database wiring must not live here.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/tasks", h.tasksCollection)
	mux.HandleFunc("/api/v1/tasks/", h.taskResource)
}

func (h *Handler) tasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTask(w, r)
	case http.MethodGet:
		h.ListTasks(w, r)
	default:
		w.Header().Set("Allow", strings.Join([]string{http.MethodPost, http.MethodGet}, ", "))
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *Handler) taskResource(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	if id == "" || strings.Contains(id, "/") {
		writeServiceError(w, ErrTaskNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetTask(w, r, id)
	case http.MethodPatch:
		h.UpdateTask(w, r, id)
	case http.MethodDelete:
		h.DeleteTask(w, r, id)
	default:
		w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPatch, http.MethodDelete}, ", "))
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid JSON request body")
		return
	}

	response, err := h.service.CreateTask(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	page, err := parseOptionalInt(r.URL.Query().Get("page"))
	if err != nil {
		writeServiceError(w, ErrInvalidPagination)
		return
	}
	limit, err := parseOptionalInt(r.URL.Query().Get("limit"))
	if err != nil {
		writeServiceError(w, ErrInvalidPagination)
		return
	}

	response, err := h.service.ListTasks(r.Context(), page, limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	response, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request, id string) {
	var req UpdateTaskRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid JSON request body")
		return
	}

	response, err := h.service.UpdateTask(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dest any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(dest)
}

func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}

	return strconv.Atoi(value)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTaskNotFound):
		writeError(w, http.StatusNotFound, "TASK_NOT_FOUND", "task not found")
	case errors.Is(err, ErrInvalidTitle):
		writeError(w, http.StatusBadRequest, "INVALID_TITLE", "invalid title")
	case errors.Is(err, ErrInvalidStatus):
		writeError(w, http.StatusBadRequest, "INVALID_STATUS", "invalid status")
	case errors.Is(err, ErrInvalidPagination):
		writeError(w, http.StatusBadRequest, "INVALID_PAGINATION", "invalid pagination")
	case errors.Is(err, ErrNoFieldsToUpdate):
		writeError(w, http.StatusBadRequest, "NO_FIELDS_TO_UPDATE", "no fields to update")
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}
