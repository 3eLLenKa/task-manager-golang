package handlers

import (
	"context"
	"todo/db/internal/domain/models"
	dbpb "todo/proto/db/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DB interface {
	CreateTask(ctx context.Context, title, description string, completed bool) error
	GetTask(ctx context.Context, id int64) (models.Task, error)
	EditTask(ctx context.Context, id int64, title, description string) error
	DeleteTask(ctx context.Context, id int64) error
	CompleteTask(ctx context.Context, id int64) error
	ListTasks(ctx context.Context) ([]models.Task, error)
	ListCompletedTasks(ctx context.Context) ([]models.Task, error)
	ListNotCompletedTasks(ctx context.Context) ([]models.Task, error)
}

type ServerApi struct {
	dbpb.UnimplementedTaskServiceServer
	db DB
}

func Register(gRPCserver *grpc.Server, db DB) {
	dbpb.RegisterTaskServiceServer(gRPCserver, &ServerApi{db: db})
}

func (s *ServerApi) CreateTask(ctx context.Context, in *dbpb.TaskRequest) (*dbpb.TaskResponse, error) {
	if in.Title == "" || in.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	if err := s.db.CreateTask(ctx,
		in.GetTitle(),
		in.GetDescription(),
		in.GetCompleted(),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbpb.TaskResponse{
		Status:  codes.OK.String(),
		Message: "created",
	}, nil
}

func (s *ServerApi) GetTask(ctx context.Context, in *dbpb.TaskId) (*dbpb.TaskItemResponse, error) {
	if in.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	data, err := s.db.GetTask(ctx, in.GetId())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	task := &dbpb.TaskItem{
		Id:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Completed:   data.Completed,
		CreatedAt:   data.CreatedAt,
		CompletedAt: data.CompletedAt,
	}

	return &dbpb.TaskItemResponse{
		Task: task,
	}, nil
}

func (s *ServerApi) EditTask(ctx context.Context, in *dbpb.EditTaskRequest) (*dbpb.TaskResponse, error) {
	if in.Id < 0 || in.Title == "" || in.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	if err := s.db.EditTask(ctx,
		in.GetId(),
		in.GetTitle(),
		in.GetDescription(),
	); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbpb.TaskResponse{
		Status:  codes.OK.String(),
		Message: "success",
	}, nil
}

func (s *ServerApi) DeleteTask(ctx context.Context, in *dbpb.TaskId) (*dbpb.TaskResponse, error) {
	if in.Id < 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	if err := s.db.DeleteTask(ctx, in.GetId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbpb.TaskResponse{
		Status:  codes.OK.String(),
		Message: "success",
	}, nil
}

func (s *ServerApi) CompleteTask(ctx context.Context, in *dbpb.TaskId) (*dbpb.TaskResponse, error) {
	if in.Id < 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	if err := s.db.CompleteTask(ctx, in.GetId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbpb.TaskResponse{
		Status:  codes.OK.String(),
		Message: "success",
	}, nil
}

func (s *ServerApi) ListTasks(ctx context.Context, in *dbpb.Empty) (*dbpb.TasksResponse, error) {
	data, err := s.db.ListTasks(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var tasks []*dbpb.TaskItem

	for _, v := range data {
		tasks = append(tasks, &dbpb.TaskItem{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Completed:   v.Completed,
			CreatedAt:   v.CreatedAt,
			CompletedAt: v.CompletedAt,
		})
	}

	return &dbpb.TasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *ServerApi) ListCompletedTasks(ctx context.Context, in *dbpb.Empty) (*dbpb.TasksResponse, error) {
	data, err := s.db.ListCompletedTasks(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var tasks []*dbpb.TaskItem

	for _, v := range data {
		tasks = append(tasks, &dbpb.TaskItem{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Completed:   v.Completed,
			CreatedAt:   v.CreatedAt,
			CompletedAt: v.CompletedAt,
		})
	}

	return &dbpb.TasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *ServerApi) ListNotCompletedTasks(ctx context.Context, in *dbpb.Empty) (*dbpb.TasksResponse, error) {
	data, err := s.db.ListNotCompletedTasks(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var tasks []*dbpb.TaskItem

	for _, v := range data {
		tasks = append(tasks, &dbpb.TaskItem{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			Completed:   v.Completed,
			CreatedAt:   v.CreatedAt,
			CompletedAt: v.CompletedAt,
		})
	}

	return &dbpb.TasksResponse{
		Tasks: tasks,
	}, nil
}
