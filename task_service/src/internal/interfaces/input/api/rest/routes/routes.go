package routes

import (
	"net/http"
	"task_service/src/internal/interfaces/input/api/rest/handler"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler *handler.TaskHandler) http.Handler {

	r := chi.NewRouter()

	r.Post("/tasks", userHandler.CreateTaskHandler)
	r.Put("/tasks/{id}", userHandler.UpdateTaskHandler)
	r.Get("/tasks", userHandler.ListTasksHandler)
	r.Patch("/tasks/{id}/complete", userHandler.CompleteTaskHandler)


	return r
}
