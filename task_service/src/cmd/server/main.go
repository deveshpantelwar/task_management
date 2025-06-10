package main

import (
	"database/sql"
	"log"
	"net/http"
	"task_service/src/internal/adaptors/persistence"
	"task_service/src/internal/interfaces/input/api/rest/handler"
	"task_service/src/internal/interfaces/input/api/rest/routes"
	"task_service/src/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL connection
	dsn := "host=localhost port=5432 user=postgres password=3695 dbname=taskmanagement sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Init repo -> usecase -> handler
	taskRepo := persistence.NewTaskRepo(db)
	taskUC := usecase.NewTaskUsecase(taskRepo)
	taskHandler := handler.NewTaskHandler(taskUC)

	// Setup router
	r := routes.InitRoutes(taskHandler)

	log.Println("Task service running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
