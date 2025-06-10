// package main

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"
// 	"task_service/src/internal/adaptors/persistence"
// 	"task_service/src/internal/adaptors/persistence/redis"
// 	"task_service/src/internal/interfaces/input/api/rest/handler"
// 	"task_service/src/internal/interfaces/input/api/rest/routes"
// 	"task_service/src/internal/usecase"

// 	_ "github.com/lib/pq"
// )

// func main() {
// 	// PostgreSQL connection
// 	dsn := "host=localhost port=5432 user=postgres password=3695 dbname=taskmanagement sslmode=disable"
// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("failed to connect to DB: %v", err)
// 	}
// 	defer db.Close()

// 	// Init repo -> usecase -> handler
// 	taskRepo := persistence.NewTaskRepo(db)
// 	taskUC := usecase.NewTaskUsecase(taskRepo)
// 	taskHandler := handler.NewTaskHandler(taskUC)

// 	publisher := redis.NewRedisPublisher("localhost:6379", "")
// 	taskService := usecase.NewTaskUsecase(taskRepo, publisher)

// 	// Setup router
// 	r := routes.InitRoutes(taskHandler)

// 	log.Println("Task service running on :8080")
// 	if err := http.ListenAndServe(":8080", r); err != nil {
// 		log.Fatalf("failed to start server: %v", err)
// 	}
// }

package main

import (
	"database/sql"
	"log"
	"net/http"
	"task_service/src/internal/adaptors/persistence"
	"task_service/src/internal/adaptors/persistence/redis"
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

	// Initialize repository
	taskRepo := persistence.NewTaskRepo(db)

	// Initialize Redis publisher
	publisher := redis.NewRedisPublisher("localhost:6379", "")

	// Initialize usecase with repo and Redis publisher
	taskUC := usecase.NewTaskUsecase(taskRepo, publisher)

	// Initialize handler
	taskHandler := handler.NewTaskHandler(taskUC)

	// Setup router
	r := routes.InitRoutes(taskHandler)

	log.Println("Task service running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
