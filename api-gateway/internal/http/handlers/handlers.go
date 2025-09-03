package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo/api/internal/domain/models"

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
	todo Todo
}

func New(todo Todo) *Handlers {
	return &Handlers{
		todo: todo,
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

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todo.ListTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) ListCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todo.ListCompletedTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) ListNotCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todo.ListNotCompletedTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tasks)
}
