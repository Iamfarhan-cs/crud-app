package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, task Task) (Task, error) {
	const query = `
		INSERT INTO tasks (
			id,
			title,
			description,
			status,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, status, created_at, updated_at, deleted_at
	`

	created, err := scanTask(r.db.QueryRowContext(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
		task.DeletedAt,
	))
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}

	return created, nil
}

func (r *PostgresRepository) FindActiveByID(ctx context.Context, id string) (Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at, deleted_at
		FROM tasks
		WHERE id = $1
			AND deleted_at IS NULL
	`

	task, err := scanTask(r.db.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("find active task by id: %w", err)
	}

	return task, nil
}

func (r *PostgresRepository) ListActive(ctx context.Context, limit int, offset int) ([]Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list active tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan active task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active tasks: %w", err)
	}

	return tasks, nil
}

func (r *PostgresRepository) UpdateActive(ctx context.Context, task Task) (Task, error) {
	const query = `
		UPDATE tasks
		SET title = $2,
			description = $3,
			status = $4,
			updated_at = $5
		WHERE id = $1
			AND deleted_at IS NULL
		RETURNING id, title, description, status, created_at, updated_at, deleted_at
	`

	updated, err := scanTask(r.db.QueryRowContext(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("update active task: %w", err)
	}

	return updated, nil
}

func (r *PostgresRepository) SoftDelete(ctx context.Context, id string) error {
	const query = `
		UPDATE tasks
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
			AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("soft delete task rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (Task, error) {
	var task Task
	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	); err != nil {
		return Task{}, err
	}

	return task, nil
}
