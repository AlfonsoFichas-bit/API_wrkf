package routes

import (
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/middleware"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// SetupRoutes configures the application routes.
func SetupRoutes(e *echo.Echo, userHandler *handlers.UserHandler, projectHandler *handlers.ProjectHandler, sprintHandler *handlers.SprintHandler, userStoryHandler *handlers.UserStoryHandler, taskHandler *handlers.TaskHandler, notificationHandler *handlers.NotificationHandler, rubricHandler *handlers.RubricHandler, evaluationHandler *handlers.EvaluationHandler, burndownHandler *handlers.BurndownHandler, websocketHandler *handlers.WebsocketHandler, jwtSecret string) {
	// --- Swagger Route ---
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// --- Public Routes ---
	e.POST("/login", userHandler.Login)
	e.POST("/register", userHandler.CreateUser)
	e.POST("/create-admin", userHandler.CreateAdminUser)

	// --- General Authenticated Routes ---
	api := e.Group("/api")
	api.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// User routes
	api.GET("/me", userHandler.GetCurrentUser)
	api.GET("/users/:id", userHandler.GetUser)

	// Notification routes
	api.GET("/notifications", notificationHandler.GetUserNotifications)
	api.POST("/notifications/read/all", notificationHandler.MarkAllAsRead)
	api.POST("/notifications/:id/read", notificationHandler.MarkAsRead)

	// Project routes
	api.POST("/projects", projectHandler.CreateProject)
	api.GET("/projects", projectHandler.GetAllProjects)
	api.GET("/projects/:id", projectHandler.GetProjectByID)
	api.PUT("/projects/:id", projectHandler.UpdateProject)
	api.DELETE("/projects/:id", projectHandler.DeleteProject)
	api.GET("/projects/:id/unassigned-users", projectHandler.GetUnassignedUsers)
	api.GET("/projects/:id/members", projectHandler.GetProjectMembers)
	api.GET("/projects/:id/board", projectHandler.GetProjectBoard)

	// Rubric routes
	api.POST("/rubrics", rubricHandler.CreateRubric)
	api.GET("/rubrics", rubricHandler.GetAllRubrics)
	api.GET("/rubrics/:id", rubricHandler.GetRubricByID)
	api.PUT("/rubrics/:id", rubricHandler.UpdateRubric)
	api.DELETE("/rubrics/:id", rubricHandler.DeleteRubric)
	api.POST("/rubrics/:id/duplicate", rubricHandler.DuplicateRubric)

	// Evaluation routes
	api.POST("/evaluations", evaluationHandler.CreateEvaluation)
	api.GET("/evaluations/:id", evaluationHandler.GetEvaluation)
	api.GET("/students/:studentId/evaluations", evaluationHandler.GetStudentEvaluations)

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
	api.PUT("/tasks/:taskId/status", taskHandler.UpdateTaskStatus)
	api.POST("/tasks/:id/comments", taskHandler.AddComment)

	// Sprint routes
	api.POST("/projects/:id/sprints", sprintHandler.CreateSprint)
	api.GET("/projects/:id/sprints", sprintHandler.GetSprintsByProjectID)
	api.GET("/sprints/:sprintId", sprintHandler.GetSprintByID)
	api.PUT("/sprints/:sprintId", sprintHandler.UpdateSprint)
	api.DELETE("/sprints/:sprintId", sprintHandler.DeleteSprint)
	api.POST("/sprints/:sprintId/userstories", userStoryHandler.AssignUserStoryToSprint)
	api.GET("/sprints/:id/burndown", burndownHandler.GetBurndownChart)

	// --- Admin-Only Routes ---
	admin := e.Group("/api/admin")
	admin.Use(middleware.JWTAuthMiddleware(jwtSecret))
	admin.Use(middleware.AdminAuthMiddleware)

	admin.GET("/users", userHandler.GetAllUsers)
	admin.POST("/users", userHandler.CreateUser)
	admin.POST("/users/admin", userHandler.CreateAdminUser)
	admin.PUT("/users/:id", userHandler.UpdateUser)
	admin.DELETE("/users/:id", userHandler.DeleteUser)

	admin.POST("/projects/:id/members", projectHandler.AddMemberToProject)

	// --- WebSocket Route (No JWT Middleware) ---
	e.GET("/api/ws/projects/:id/board", websocketHandler.ServeWs)
}
