package handlers

import (
	"net/http"
	"todo/api/internal/domain/models"
)

type Todo interface {
	CreateTask(name, description string) error
	GetTask(id int) (models.Task, error)
	EditTask(id int, name, description string) error
	DeleteTask(id int) error
	CompleteTask(id int) error
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

}

func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) EditTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) ListTasksHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) ListCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) ListNotCompletedTasksHandler(w http.ResponseWriter, r *http.Request) {

}
