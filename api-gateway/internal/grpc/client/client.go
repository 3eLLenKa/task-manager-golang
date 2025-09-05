package client

import (
	"context"
	"fmt"
	"time"
	"todo/api/internal/domain/models"
	dbpb "todo/proto/db/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client dbpb.TaskServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	newClient := dbpb.NewTaskServiceClient(conn)

	return &Client{
		client: newClient,
	}, nil
}

func (c *Client) CreateTask(name, description string) error {
	const op = "client.CreateTask"

	_, err := c.client.CreateTask(context.Background(), &dbpb.TaskRequest{
		Title:       name,
		Description: description,
		Completed:   false,
	})

	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (c *Client) GetTask(id int64) (models.Task, error) {
	const op = "client.GetTask"

	task, err := c.client.GetTask(context.Background(), &dbpb.TaskId{
		Id: id,
	})

	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %v", op, err)
	}

	model := task.Task

	var completedAt *time.Time

	if model.CompletedAt != nil {
		t := model.CompletedAt.AsTime()
		completedAt = &t
	}

	return models.Task{
		Id:          model.Id,
		Name:        model.Title,
		Description: model.Description,
		Completed:   model.Completed,
		CreatedAt:   model.CreatedAt.AsTime(),
		CompletedAt: completedAt,
	}, nil
}

func (c *Client) EditTask(id int64, name, description string) error {
	const op = "client.EditTask"

	_, err := c.client.EditTask(context.Background(), &dbpb.EditTaskRequest{
		Id:          id,
		Title:       name,
		Description: description,
	})

	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (c *Client) DeleteTask(id int64) error {
	const op = "client.DeleteTask"

	_, err := c.client.DeleteTask(context.Background(), &dbpb.TaskId{
		Id: id,
	})

	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (c *Client) CompleteTask(id int64) error {
	const op = "client.CompleteTask"

	_, err := c.client.CompleteTask(context.Background(), &dbpb.TaskId{
		Id: id,
	})

	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (c *Client) ListTasks() ([]models.Task, error) {
	const op = "client.ListTasks"

	tasks, err := c.client.ListTasks(context.Background(), &dbpb.Empty{})

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var resp []models.Task

	for _, v := range tasks.Tasks {
		var completedAt *time.Time

		if v.CompletedAt != nil {
			t := v.CompletedAt.AsTime()
			completedAt = &t
		}

		resp = append(resp, models.Task{
			Id:          v.Id,
			Name:        v.Title,
			Description: v.Description,
			Completed:   v.Completed,
			CreatedAt:   v.CreatedAt.AsTime(),
			CompletedAt: completedAt,
		})
	}

	return resp, nil
}

func (c *Client) ListCompletedTasks() ([]models.Task, error) {
	const op = "client.ListCompletedTasks"

	tasks, err := c.client.ListCompletedTasks(context.Background(), &dbpb.Empty{})

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var resp []models.Task

	for _, v := range tasks.Tasks {
		var completedAt *time.Time

		if v.CompletedAt != nil {
			t := v.CompletedAt.AsTime()
			completedAt = &t
		}

		resp = append(resp, models.Task{
			Id:          v.Id,
			Name:        v.Title,
			Description: v.Description,
			Completed:   v.Completed,
			CreatedAt:   v.CreatedAt.AsTime(),
			CompletedAt: completedAt,
		})
	}

	return resp, nil
}
func (c *Client) ListNotCompletedTasks() ([]models.Task, error) {
	const op = "client.ListNotCompletedTasks"

	tasks, err := c.client.ListNotCompletedTasks(context.Background(), &dbpb.Empty{})

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var resp []models.Task

	for _, v := range tasks.Tasks {
		var completedAt *time.Time

		if v.CompletedAt != nil {
			t := v.CompletedAt.AsTime()
			completedAt = &t
		}

		resp = append(resp, models.Task{
			Id:          v.Id,
			Name:        v.Title,
			Description: v.Description,
			Completed:   v.Completed,
			CreatedAt:   v.CreatedAt.AsTime(),
			CompletedAt: completedAt,
		})
	}

	return resp, nil
}
