package main

import (
	"database/sql"
	"log"
	"net/http"
	"task_management/task_service/src/internal/adaptors/external"

	"task_management/task_service/src/internal/adaptors/persistence/redis"
	persistence "task_management/task_service/src/internal/adaptors/persistence/task_repo"
	"task_management/task_service/src/internal/interfaces/input/api/rest/handler"
	"task_management/task_service/src/internal/interfaces/input/api/rest/middleware"
	"task_management/task_service/src/internal/interfaces/input/api/rest/routes"
	"task_management/task_service/src/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=3695 dbname=taskmanagement sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	taskRepo := persistence.NewTaskRepo(db)
	publisher := redis.NewRedisPublisher("localhost:6379", "")
	taskUC := usecase.NewTaskUsecase(taskRepo, publisher)
	taskHandler := handler.NewTaskHandler(taskUC)

	userClient := external.NewUserServiceClient("localhost:50051")
	authMiddleware := middleware.NewAuthMiddleware(userClient)

	r := routes.InitRoutes(taskHandler, authMiddleware)

	log.Println("Task service running on :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
