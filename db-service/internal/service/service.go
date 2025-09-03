package service

import (
	"context"
	"fmt"
	"log/slog"
	"todo/db/internal/domain/models"
	"todo/db/internal/lib/sl"
)

type TaskProvider interface {
	Save(ctx context.Context, title, description string) (int64, error)
	Get(ctx context.Context, id int64) (models.Task, error)
	Update(ctx context.Context, id int64, title, description string) error
	Remove(ctx context.Context, id int64) error
	Complete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.Task, error)
	ListCompleted(ctx context.Context) ([]models.Task, error)
	ListNotCompleted(ctx context.Context) ([]models.Task, error)
}

type TaskCache interface {
	SetTask(ctx context.Context, task models.Task) error
	GetTask(ctx context.Context, id int64) (models.Task, error)
	DelTask(ctx context.Context, id int64) error
}

type TaskService struct {
	log          *slog.Logger
	taskProvider TaskProvider
	taskCache    TaskCache
}

func New(
	log *slog.Logger,
	taskProvider TaskProvider,
	taskCache TaskCache,
) *TaskService {
	return &TaskService{
		log:          log,
		taskProvider: taskProvider,
		taskCache:    taskCache,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title, description string, completed bool) error {
	const op = "service.CreateTask"

	log := s.log.With(
		slog.String("op", op),
	)

	id, err := s.taskProvider.Save(ctx, title, description)

	if err != nil {
		log.Error("task not created", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	task, err := s.taskProvider.Get(ctx, id)

	if err == nil {
		_ = s.taskCache.SetTask(ctx, task)
	}

	return nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (models.Task, error) {
	const op = "service.GetTask"

	log := s.log.With(
		slog.String("op", op),
	)

	task, err := s.taskProvider.Get(ctx, id)

	if err != nil {
		log.Error("task not found", sl.Err(err))
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	_ = s.taskCache.SetTask(ctx, task)

	return task, nil
}

func (s *TaskService) EditTask(ctx context.Context, id int64, title, description string) error {
	const op = "service.EditTask"

	log := s.log.With(
		slog.String("op", op),
	)

	if err := s.taskProvider.Update(ctx, id, title, description); err != nil {
		log.Error("task not updated", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = s.taskCache.DelTask(ctx, id)

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	const op = "service.DeleteTask"

	log := s.log.With(
		slog.String("op", op),
	)

	if err := s.taskProvider.Remove(ctx, id); err != nil {
		log.Error("task not deleted", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = s.taskCache.DelTask(ctx, id)

	return nil
}

func (s *TaskService) CompleteTask(ctx context.Context, id int64) error {
	const op = "service.CompleteTask"

	log := s.log.With(
		slog.String("op", op),
	)

	if err := s.taskProvider.Complete(ctx, id); err != nil {
		log.Error("task not completed", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = s.taskCache.DelTask(ctx, id)

	return nil
}

func (s *TaskService) ListTasks(ctx context.Context) ([]models.Task, error) {
	const op = "service.ListTasks"

	log := s.log.With(
		slog.String("op", op),
	)

	tasks, err := s.taskProvider.List(ctx)

	if err != nil {
		log.Error("internal error", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *TaskService) ListCompletedTasks(ctx context.Context) ([]models.Task, error) {
	const op = "service.ListCompletedTasks"

	log := s.log.With(
		slog.String("op", op),
	)

	tasks, err := s.taskProvider.ListCompleted(ctx)

	if err != nil {
		log.Error("internal error", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *TaskService) ListNotCompletedTasks(ctx context.Context) ([]models.Task, error) {
	const op = "service.ListNotCompletedTasks"

	log := s.log.With(
		slog.String("op", op),
	)

	tasks, err := s.taskProvider.ListNotCompleted(ctx)

	if err != nil {
		log.Error("internal error", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}
