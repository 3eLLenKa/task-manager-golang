package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"todo/db/internal/domain/models"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("redis: cache miss")
	ErrInternal  = errors.New("redis: internal error")
)

type RedisStorage struct {
	client *redis.Client
	ttl    time.Duration
}

func New(dsn string, ttl time.Duration) (*RedisStorage, error) {
	const op = "internal.storage.redis.New"

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &RedisStorage{
		client: client,
		ttl:    ttl,
	}, nil
}

func (s *RedisStorage) SetTask(ctx context.Context, task models.Task) error {
	key := fmt.Sprintf("task:%d", task.ID)

	data, err := json.Marshal(task)

	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, s.ttl).Err()
}

func (s *RedisStorage) GetTask(ctx context.Context, id int64) (models.Task, error) {
	key := fmt.Sprintf("task:%d", id)

	data, err := s.client.Get(ctx, key).Result()

	if err != nil {
		return models.Task{}, err
	}

	var task models.Task

	if err := json.Unmarshal([]byte(data), &task); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *RedisStorage) DelTask(ctx context.Context, id int64) error {
	key := fmt.Sprintf("task:%d", id)

	return s.client.Del(ctx, key).Err()
}
