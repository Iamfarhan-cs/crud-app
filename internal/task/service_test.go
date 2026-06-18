package task

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"
)

type fakeRepository struct {
	tasks     map[string]Task
	createErr error
	listErr   error
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{tasks: make(map[string]Task)}
}

func (r *fakeRepository) Create(ctx context.Context, task Task) (Task, error) {
	if r.createErr != nil {
		return Task{}, r.createErr
	}
	r.tasks[task.ID] = task
	return task, nil
}

func (r *fakeRepository) FindActiveByID(ctx context.Context, id string) (Task, error) {
	task, ok := r.tasks[id]
	if !ok || task.DeletedAt != nil {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (r *fakeRepository) ListActive(ctx context.Context, limit int, offset int) ([]Task, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}

	tasks := make([]Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		if task.DeletedAt == nil {
			tasks = append(tasks, task)
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	if offset >= len(tasks) {
		return []Task{}, nil
	}

	end := offset + limit
	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[offset:end], nil
}

func (r *fakeRepository) UpdateActive(ctx context.Context, task Task) (Task, error) {
	current, ok := r.tasks[task.ID]
	if !ok || current.DeletedAt != nil {
		return Task{}, ErrTaskNotFound
	}
	r.tasks[task.ID] = task
	return task, nil
}

func (r *fakeRepository) SoftDelete(ctx context.Context, id string) error {
	task, ok := r.tasks[id]
	if !ok || task.DeletedAt != nil {
		return ErrTaskNotFound
	}

	now := time.Now().UTC()
	task.DeletedAt = &now
	task.UpdatedAt = now
	r.tasks[id] = task

	return nil
}

func TestServiceCreateTaskSuccessGeneratesIDTrimsTitleAndDefaultsStatus(t *testing.T) {
	service := NewService(newFakeRepository())

	got, err := service.CreateTask(context.Background(), CreateTaskRequest{
		Title: "  Write tests  ",
	})
	if err != nil {
		t.Fatalf("CreateTask returned error: %v", err)
	}

	if got.ID == "" {
		t.Fatal("expected generated ID")
	}
	if got.Title != "Write tests" {
		t.Fatalf("expected trimmed title, got %q", got.Title)
	}
	if got.Status != StatusPending {
		t.Fatalf("expected default status %q, got %q", StatusPending, got.Status)
	}
}

func TestServiceCreateTaskRejectsWhitespaceOnlyTitle(t *testing.T) {
	service := NewService(newFakeRepository())

	_, err := service.CreateTask(context.Background(), CreateTaskRequest{Title: "   "})
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestServiceCreateTaskRejectsInvalidStatus(t *testing.T) {
	service := NewService(newFakeRepository())
	status := Status("blocked")

	_, err := service.CreateTask(context.Background(), CreateTaskRequest{
		Title:  "Write tests",
		Status: &status,
	})
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestServiceUpdateTaskSuccessUpdatesProvidedFieldsOnly(t *testing.T) {
	repo := newFakeRepository()
	createdAt := time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC)
	description := "original description"
	original := Task{
		ID:          "task-1",
		Title:       "Original title",
		Description: &description,
		Status:      StatusPending,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	repo.tasks[original.ID] = original
	service := NewService(repo)

	title := "  Updated title  "
	got, err := service.UpdateTask(context.Background(), original.ID, UpdateTaskRequest{
		Title: &title,
	})
	if err != nil {
		t.Fatalf("UpdateTask returned error: %v", err)
	}

	if got.Title != "Updated title" {
		t.Fatalf("expected updated title, got %q", got.Title)
	}
	if got.Description == nil || *got.Description != description {
		t.Fatalf("expected description to stay %q, got %#v", description, got.Description)
	}
	if got.Status != StatusPending {
		t.Fatalf("expected status to stay %q, got %q", StatusPending, got.Status)
	}
	if !got.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected created_at to stay %s, got %s", createdAt, got.CreatedAt)
	}
	if !got.UpdatedAt.After(createdAt) {
		t.Fatalf("expected updated_at to move forward, got %s", got.UpdatedAt)
	}
}

func TestServiceUpdateTaskRejectsNoFields(t *testing.T) {
	service := NewService(newFakeRepository())

	_, err := service.UpdateTask(context.Background(), "task-1", UpdateTaskRequest{})
	if !errors.Is(err, ErrNoFieldsToUpdate) {
		t.Fatalf("expected ErrNoFieldsToUpdate, got %v", err)
	}
}

func TestServiceDeleteTaskHidesTaskFromFutureGetTask(t *testing.T) {
	repo := newFakeRepository()
	now := time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC)
	repo.tasks["task-1"] = Task{
		ID:        "task-1",
		Title:     "Delete me",
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	service := NewService(repo)

	if err := service.DeleteTask(context.Background(), "task-1"); err != nil {
		t.Fatalf("DeleteTask returned error: %v", err)
	}

	_, err := service.GetTask(context.Background(), "task-1")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestServiceListTasksRejectsLimitAboveMaxListLimit(t *testing.T) {
	service := NewService(newFakeRepository())

	_, err := service.ListTasks(context.Background(), 1, MaxListLimit+1)
	if !errors.Is(err, ErrInvalidPagination) {
		t.Fatalf("expected ErrInvalidPagination, got %v", err)
	}
}
