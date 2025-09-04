package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo/api/internal/domain/models"
	producer "todo/api/internal/kafka"

	"github.com/go-chi/chi/v5"
)

type Todo interface {
	CreateTask(name, description string) error
	GetTask(id int64) (models.Task, error)
	EditTask(id int64, name, description string) error
	DeleteTask(id int64) error
	CompleteTask(id int64) error
	ListTasks() ([]models.Task, error)
	ListCompletedTasks() ([]models.Task, error)
	ListNotCompletedTasks() ([]models.Task, error)
}

type Handlers struct {
	producer *producer.Producer
	todo     Todo
}

func New(todo Todo, producer *producer.Producer) *Handlers {
	return &Handlers{
		todo:     todo,
		producer: producer,
	}
}

func (h *Handlers) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.todo.CreateTask(req.Name, req.Description); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=create_task name=%s",
			time.Now().Format(time.RFC3339), req.Name),
	)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.todo.GetTask(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=get_task id=%d",
			time.Now().Format(time.RFC3339), id),
	)

	_ = json.NewEncoder(w).Encode(task)
}

func (h *Handlers) EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.todo.EditTask(id, req.Name, req.Description); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=edit_task id=%d name=%s",
			time.Now().Format(time.RFC3339), id, req.Name),
	)

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.todo.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=delete_task id=%d",
			time.Now().Format(time.RFC3339), id),
	)

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.todo.CompleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=complete_task id=%d",
			time.Now().Format(time.RFC3339), id),
	)

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todo.ListTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=list_tasks count=%d",
			time.Now().Format(time.RFC3339), len(tasks)),
	)

	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) ListCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todo.ListCompletedTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=list_completed_tasks count=%d",
			time.Now().Format(time.RFC3339), len(tasks)),
	)

	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) ListNotCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todo.ListNotCompletedTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = h.producer.Publish(
		fmt.Sprintf("time=%s action=list_not_completed_tasks count=%d",
			time.Now().Format(time.RFC3339), len(tasks)),
	)

	_ = json.NewEncoder(w).Encode(tasks)
}
