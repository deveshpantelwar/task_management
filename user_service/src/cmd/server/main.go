package main

import (
	"fmt"
	"log"
	"net/http"
	"task-management/user-service/src/internal/adaptors/persistance"
	"task-management/user-service/src/internal/config"
	"task-management/user-service/src/internal/interfaces/input/api/rest/handler"
	"task-management/user-service/src/internal/interfaces/input/api/rest/routes"
	"task-management/user-service/src/internal/usecase"
)

const (
	port = ":8080"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	database, err := persistance.ConnectToDatabase(config)
	if err != nil {
		log.Fatalf("could not connect to database :- %v", err)
	}

	userRepository := persistance.NewUserRepo(database)
	userUsecase := usecase.NewUserService(userRepository, config.JWT_SECRET)
	userHandler := handler.NewUserHandler(userUsecase)

	

	router := routes.InitRoutes(userHandler)

	fmt.Printf("Starting server on port 8080 \n")

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Failed to connect to server : %v", err)
	}
}
