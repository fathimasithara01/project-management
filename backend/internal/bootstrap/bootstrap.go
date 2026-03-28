package bootstrap

import (
	"log"
	"time"

	"project-management/config"
	"project-management/internal/handler"
	"project-management/internal/repositories"
	"project-management/internal/usecase"
	"project-management/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func NewRouter() (*gin.Engine, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	userUC := usecase.NewUserService(userRepo)
	projectUC := usecase.NewProjectService(projectRepo)
	taskUC := usecase.NewTaskService(taskRepo)

	authHandler := handler.NewAuthHandler(userUC)
	userHandler := handler.NewUserHandler(userUC)
	projectHandler := handler.NewProjectHandler(projectUC)
	taskHandler := handler.NewTaskHandler(taskUC)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(router, authHandler, userHandler, projectHandler, taskHandler)

	return router, nil
}
