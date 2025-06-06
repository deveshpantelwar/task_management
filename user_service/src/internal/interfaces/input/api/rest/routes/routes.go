package routes

import (
	"net/http"
	"task-management/user-service/src/internal/interfaces/input/api/rest/handler"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler *handler.UserHandler) http.Handler {

	r := chi.NewRouter()

	r.Post("/register", userHandler.RegisterHandler)
	r.Post("/login", userHandler.LoginHandler)

	r.Get("/profile", userHandler.GetUserProfileHandler)
	r.Put("/update", userHandler.UpdateUserHandler)

	return r
}
