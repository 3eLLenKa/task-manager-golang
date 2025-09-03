package router

import (
	"net/http"
	"todo/api/internal/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	handlers *handlers.Handlers
}

func New(handlers *handlers.Handlers) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) InitRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Route("/api/v1/todos", func(ch chi.Router) {
		ch.Get("/", r.handlers.ListTasksHandler)   // GET /api/v1/todos
		ch.Post("/", r.handlers.CreateTaskHandler) // POST /api/v1/todos

		ch.Get("/{id}", r.handlers.GetTaskHandler)                 // GET /api/v1/todos/{id}
		ch.Put("/{id}", r.handlers.EditTaskHandler)                // PUT /api/v1/todos/{id}
		ch.Delete("/{id}", r.handlers.DeleteTaskHandler)           // DELETE /api/v1/todos/{id}
		ch.Patch("/{id}/complete", r.handlers.CompleteTaskHandler) // PATCH /api/v1/todos/{id}/complete

		ch.Get("/completed", r.handlers.ListCompletedTasksHandler)  // GET /api/v1/todos/completed
		ch.Get("/pending", r.handlers.ListNotCompletedTasksHandler) // GET /api/v1/todos/pending
	})

	return router
}
