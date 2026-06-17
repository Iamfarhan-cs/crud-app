package task

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultListLimit = 20
	MaxListLimit     = 100
)

// Service owns task business rules and use-case orchestration.
// It should validate domain behavior and coordinate repository calls.
// Service-level security work includes enforcing business validation that must
// hold even when requests come from different HTTP handlers or future clients.
// Future performance work should pass request contexts through to repository
// calls so timeouts and cancellations are respected.
// Future concurrency rules such as stale update conflicts and idempotent create behavior belong here.
// HTTP request parsing, response formatting, SQL, authentication, authorization,
// and database connection setup must not live here.
type Service struct {
	repo Repository
}

// NewService prepares the task service with its persistence dependency.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(ctx context.Context, req CreateTaskRequest) (TaskResponse, error) {
	title, err := validateTitle(req.Title)
	if err != nil {
		return TaskResponse{}, err
	}

	status := StatusPending
	if req.Status != nil {
		status = *req.Status
	}
	if err := validateStatus(status); err != nil {
		return TaskResponse{}, err
	}

	now := time.Now().UTC()
	task := Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: req.Description,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := s.repo.Create(ctx, task)
	if err != nil {
		return TaskResponse{}, err
	}

	return toResponse(created), nil
}

func (s *Service) GetTask(ctx context.Context, id string) (TaskResponse, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return TaskResponse{}, ErrTaskNotFound
	}

	task, err := s.repo.FindActiveByID(ctx, id)
	if err != nil {
		return TaskResponse{}, err
	}

	return toResponse(task), nil
}

func (s *Service) ListTasks(ctx context.Context, page int, limit int) ([]TaskResponse, error) {
	limit, offset, err := normalizePagination(page, limit)
	if err != nil {
		return nil, err
	}

	tasks, err := s.repo.ListActive(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		responses = append(responses, toResponse(task))
	}

	return responses, nil
}

func (s *Service) UpdateTask(ctx context.Context, id string, req UpdateTaskRequest) (TaskResponse, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return TaskResponse{}, ErrTaskNotFound
	}
	if req.Title == nil && req.Description == nil && req.Status == nil {
		return TaskResponse{}, ErrNoFieldsToUpdate
	}

	task, err := s.repo.FindActiveByID(ctx, id)
	if err != nil {
		return TaskResponse{}, err
	}

	if req.Title != nil {
		title, err := validateTitle(*req.Title)
		if err != nil {
			return TaskResponse{}, err
		}
		task.Title = title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.Status != nil {
		if err := validateStatus(*req.Status); err != nil {
			return TaskResponse{}, err
		}
		task.Status = *req.Status
	}
	task.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.UpdateActive(ctx, task)
	if err != nil {
		return TaskResponse{}, err
	}

	return toResponse(updated), nil
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrTaskNotFound
	}

	return s.repo.SoftDelete(ctx, id)
}

func validateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", ErrInvalidTitle
	}

	return title, nil
}

func validateStatus(status Status) error {
	switch status {
	case StatusPending, StatusInProgress, StatusDone:
		return nil
	default:
		return ErrInvalidStatus
	}
}

func normalizePagination(page int, limit int) (int, int, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = DefaultListLimit
	}
	if page < 1 || limit < 1 || limit > MaxListLimit {
		return 0, 0, ErrInvalidPagination
	}

	return limit, (page - 1) * limit, nil
}

func toResponse(t Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
