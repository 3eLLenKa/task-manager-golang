package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo/api/internal/domain/models"
	"todo/api/internal/http/handlers"
	producer "todo/api/internal/kafka"
)

type fakeProducer struct {
	writer   *producer.Producer
	messages []string
}

func (f *fakeProducer) Publish(event string) error {
	f.messages = append(f.messages, event)
	return nil
}

func (f *fakeProducer) Close() error { return nil }

type fakeTodo struct{}

func (f *fakeTodo) CreateTask(name, description string) error { return nil }

func (f *fakeTodo) GetTask(id int64) (models.Task, error) {
	return models.Task{Id: id, Name: "Task", Description: "Desc", CreatedAt: time.Now()}, nil
}

func (f *fakeTodo) EditTask(id int64, name, description string) error { return nil }

func (f *fakeTodo) DeleteTask(id int64) error { return nil }

func (f *fakeTodo) CompleteTask(id int64) error { return nil }

func (f *fakeTodo) ListTasks() ([]models.Task, error) {
	return []models.Task{{Id: 1, Name: "Task"}}, nil
}

func (f *fakeTodo) ListCompletedTasks() ([]models.Task, error) {
	return []models.Task{{Id: 1, Name: "Task"}}, nil
}

func (f *fakeTodo) ListNotCompletedTasks() ([]models.Task, error) {
	return []models.Task{{Id: 2, Name: "Task2"}}, nil
}

func TestHandlers(t *testing.T) {
	todo := &fakeTodo{}
	prod := &fakeProducer{
		writer: nil,
	}

	h := handlers.New(todo, prod.writer)

	// CreateTask
	{
		body := map[string]string{"name": "Test", "description": "Desc"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
		w := httptest.NewRecorder()

		h.CreateTaskHandler(w, req)

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("CreateTaskHandler: ожидался 201, получили %d", w.Result().StatusCode)
		}

		if len(prod.messages) != 1 {
			t.Fatalf("CreateTaskHandler: ожидалось 1 сообщение, получили %d", len(prod.messages))
		}
	}

	// GetTask
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w := httptest.NewRecorder()

		h.GetTaskHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("GetTaskHandler: ожидался 200, получили %d", w.Result().StatusCode)
		}
	}

	// EditTask
	{
		body := map[string]string{"name": "Edited", "description": "Desc2"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(b))

		w := httptest.NewRecorder()

		h.EditTaskHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("EditTaskHandler: ожидался 200, получили %d", w.Result().StatusCode)
		}
	}

	// DeleteTask
	{
		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		w := httptest.NewRecorder()

		h.DeleteTaskHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("DeleteTaskHandler: ожидался 200, получили %d", w.Result().StatusCode)
		}
	}

	// CompleteTask
	{
		req := httptest.NewRequest(http.MethodPost, "/tasks/1/complete", nil)
		w := httptest.NewRecorder()
		h.CompleteTaskHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("CompleteTaskHandler: ожидался 200, получили %d", w.Result().StatusCode)
		}
	}

	// ListTasks
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		h.ListTasksHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ListTasksHandler: ожидался 200, получили %d", w.Result().StatusCode)
		}
	}

	// ListCompletedTasks
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks/completed", nil)
		w := httptest.NewRecorder()

		h.ListCompletedTasksHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ListCompletedTasksHandler: ожидался 200, получили %d", w.Result().StatusCode)
		}
	}

	// ListNotCompletedTasks
	{
		req := httptest.NewRequest(http.MethodGet, "/tasks/not-completed", nil)
		w := httptest.NewRecorder()

		h.ListNotCompletedTasksHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ListNotCompletedTasksHandler: ожидался 200, получили %d", w.Result().StatusCode)
		}
	}
}
