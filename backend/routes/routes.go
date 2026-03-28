package routes

import (
	"project-management/internal/handler"
	"project-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	projectHandler *handler.ProjectHandler,
	taskHandler *handler.TaskHandler,

) {
	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	users := protected.Group("/users")
	users.Use(middleware.RoleMiddleware("admin"))
	{
		users.GET("/", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("/", userHandler.CreateUser)
	}

	projects := protected.Group("/projects")
	projects.Use(middleware.RoleMiddleware("admin"))
	{
		projects.GET("/", projectHandler.ListProjects)
		projects.GET("/:id", projectHandler.GetProjectByID)
		projects.POST("/", projectHandler.CreateProject)
		projects.PUT("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)
	}

	tasks := protected.Group("/tasks")
	{
		tasks.GET("/", middleware.RoleMiddleware("admin", "developer"), taskHandler.ListTasks)
		tasks.GET("/:id", middleware.RoleMiddleware("admin", "developer"), taskHandler.GetTaskByID)
		tasks.PATCH("/:id/status", middleware.RoleMiddleware("admin", "developer"), taskHandler.UpdateTaskStatus)

		tasks.POST("/", middleware.RoleMiddleware("admin"), taskHandler.CreateTask)
		tasks.PUT("/:id", middleware.RoleMiddleware("admin"), taskHandler.UpdateTask)
		tasks.DELETE("/:id", middleware.RoleMiddleware("admin"), taskHandler.DeleteTask)
	}

}
