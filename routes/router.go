package routes

import (
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/middleware"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger" // Import echo-swagger
)

// SetupRoutes configures the application routes.
func SetupRoutes(e *echo.Echo, userHandler *handlers.UserHandler, projectHandler *handlers.ProjectHandler, sprintHandler *handlers.SprintHandler, userStoryHandler *handlers.UserStoryHandler, taskHandler *handlers.TaskHandler, jwtSecret string) {
	// --- Swagger Route ---
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// --- Public Routes ---
	e.POST("/login", userHandler.Login)

	// --- General Authenticated Routes ---
	api := e.Group("/api")
	api.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// User routes
	api.GET("/users/:id", userHandler.GetUser)

	// Project routes
	api.POST("/projects", projectHandler.CreateProject)
	api.GET("/projects", projectHandler.GetAllProjects)
	api.GET("/projects/:id", projectHandler.GetProjectByID)
	api.PUT("/projects/:id", projectHandler.UpdateProject)
	api.DELETE("/projects/:id", projectHandler.DeleteProject)

	// User Story routes
	api.POST("/projects/:id/userstories", userStoryHandler.CreateUserStory)
	api.GET("/projects/:id/userstories", userStoryHandler.GetUserStoriesByProjectID)
	api.GET("/userstories/:storyId", userStoryHandler.GetUserStoryByID)
	api.PUT("/userstories/:storyId", userStoryHandler.UpdateUserStory)
	api.DELETE("/userstories/:storyId", userStoryHandler.DeleteUserStory)

	// Task routes
	api.POST("/userstories/:storyId/tasks", taskHandler.CreateTask)
	api.GET("/userstories/:storyId/tasks", taskHandler.GetTasksByUserStoryID)
	api.PUT("/tasks/:taskId", taskHandler.UpdateTask)
	api.DELETE("/tasks/:taskId", taskHandler.DeleteTask)
	api.PUT("/tasks/:taskId/assign", taskHandler.AssignTask)
	api.PUT("/tasks/:taskId/status", taskHandler.UpdateTaskStatus) // <-- NEW

	// Sprint routes
	api.POST("/projects/:id/sprints", sprintHandler.CreateSprint)
	api.GET("/projects/:id/sprints", sprintHandler.GetSprintsByProjectID)
	api.GET("/sprints/:sprintId", sprintHandler.GetSprintByID)
	api.PUT("/sprints/:sprintId", sprintHandler.UpdateSprint)
	api.DELETE("/sprints/:sprintId", sprintHandler.DeleteSprint)
	api.POST("/sprints/:sprintId/userstories", userStoryHandler.AssignUserStoryToSprint)

	// --- Admin-Only Routes ---
	admin := e.Group("/api/admin")
	admin.Use(middleware.JWTAuthMiddleware(jwtSecret))
	admin.Use(middleware.AdminAuthMiddleware)

	// Admin user management
	admin.POST("/users", userHandler.CreateUser)
	admin.POST("/users/admin", userHandler.CreateAdminUser)

	// Admin project management
	admin.POST("/projects/:id/members", projectHandler.AddMemberToProject)
}
