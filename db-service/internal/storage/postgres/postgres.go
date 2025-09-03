package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"todo/db/internal/domain/models"
)

var (
	ErrNotFound = errors.New("postgres: not found")
	ErrConflict = errors.New("postgres: conflict")
	ErrInternal = errors.New("postgres: internal error")
)

type PGStorage struct {
	db *sql.DB
}

func New(dsn string) (*PGStorage, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	return &PGStorage{db: db}, nil
}

func (s *PGStorage) Save(ctx context.Context, title, description string) (int64, error) {
	query := `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2) RETURNING id
	`
	var id int64

	if err := s.db.QueryRowContext(ctx, query, title, description).Scan(&id); err != nil {
		return -1, ErrInternal
	}

	return id, nil
}

func (s *PGStorage) Get(ctx context.Context, id int64) (models.Task, error) {
	query := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE id = $1
	`
	var task models.Task

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Task{}, ErrNotFound
		}

		return models.Task{}, ErrInternal
	}

	return task, nil
}

func (s *PGStorage) Update(ctx context.Context, id int64, title, description string) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2
		WHERE id = $3
	`

	res, err := s.db.ExecContext(ctx, query, title, description, id)

	if err != nil {
		return ErrInternal
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return ErrInternal
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PGStorage) Remove(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return ErrInternal
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return ErrInternal
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PGStorage) Complete(ctx context.Context, id int64) error {
	query := `
		UPDATE tasks
		SET completed = true, completed_at = NOW()
		WHERE id = $1
	`

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return ErrInternal
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return ErrInternal
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PGStorage) List(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
	`

	return s.fetchTasks(ctx, query)
}

func (s *PGStorage) ListCompleted(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE completed = true
	`

	return s.fetchTasks(ctx, query)
}

func (s *PGStorage) ListNotCompleted(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT id, title, description, completed, created_at, completed_at
		FROM tasks
		WHERE completed = false
	`

	return s.fetchTasks(ctx, query)
}

func (s *PGStorage) fetchTasks(ctx context.Context, query string, args ...interface{}) ([]models.Task, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, ErrInternal
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
		); err != nil {
			return nil, ErrInternal
		}

		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil, ErrNotFound
	}

	return tasks, nil
}
