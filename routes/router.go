package routes

import (
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/middleware"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger" // Import echo-swagger
)

// SetupRoutes configures the application routes.
func SetupRoutes(e *echo.Echo, userHandler *handlers.UserHandler, projectHandler *handlers.ProjectHandler, projectService *services.ProjectService, sprintHandler *handlers.SprintHandler, userStoryHandler *handlers.UserStoryHandler, taskHandler *handlers.TaskHandler, notificationHandler *handlers.NotificationHandler, rubricHandler *handlers.RubricHandler, backlogHandler *handlers.BacklogHandler, taskBoardHandler *handlers.TaskBoardHandler, evaluationHandler *handlers.EvaluationHandler, metricHandler *handlers.MetricHandler, adminHandler *handlers.AdminHandler, reportingHandler *handlers.ReportingHandler, jwtSecret string) {
	// --- Swagger Route ---
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// --- Public Routes ---
	e.POST("/login", userHandler.Login)

	// --- General Authenticated Routes ---
	api := e.Group("/api")
	api.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// User routes
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
	api.GET("/projects/:id/members", projectHandler.GetProjectMembers, middleware.ProjectRoleAuth(projectService, models.RoleScrumMaster, models.RoleProductOwner, models.RoleTeamDeveloper))

	// Backlog routes
	api.GET("/projects/:projectId/backlog", backlogHandler.GetProductBacklog)

	// Rubric routes
	api.POST("/rubrics", rubricHandler.CreateRubric)
	api.GET("/rubrics", rubricHandler.GetAllRubrics)
	api.GET("/rubrics/:id", rubricHandler.GetRubricByID)
	api.PUT("/rubrics/:id", rubricHandler.UpdateRubric)
	api.DELETE("/rubrics/:id", rubricHandler.DeleteRubric)
	api.POST("/rubrics/:id/duplicate", rubricHandler.DuplicateRubric)

	// User Story routes
	api.POST("/projects/:id/userstories", userStoryHandler.CreateUserStory)
	api.GET("/projects/:id/userstories", userStoryHandler.GetUserStoriesByProjectID)
	api.GET("/userstories/:storyId", userStoryHandler.GetUserStoryByID)
	api.PUT("/userstories/:storyId", userStoryHandler.UpdateUserStory)
	api.DELETE("/userstories/:storyId", userStoryHandler.DeleteUserStory)
	api.PUT("/userstories/:storyId/status", backlogHandler.UpdateUserStoryStatus)

	// Task routes
	api.POST("/userstories/:storyId/tasks", taskHandler.CreateTask)
	api.GET("/userstories/:storyId/tasks", taskHandler.GetTasksByUserStoryID)
	api.PUT("/tasks/:taskId", taskHandler.UpdateTask)
	api.DELETE("/tasks/:taskId", taskHandler.DeleteTask)
	api.PUT("/tasks/:taskId/assign", taskHandler.AssignTask)
	api.PUT("/tasks/:taskId/status", taskHandler.UpdateTaskStatus) // <-- NEW
	api.POST("/tasks/:id/comments", taskHandler.AddComment)      // <-- NEW

	// Sprint routes
	api.POST("/projects/:id/sprints", sprintHandler.CreateSprint, middleware.ProjectRoleAuth(projectService, models.RoleScrumMaster, models.RoleProductOwner))
	api.GET("/projects/:id/sprints", sprintHandler.GetSprintsByProjectID)
	api.GET("/sprints/:sprintId", sprintHandler.GetSprintByID)
	api.PUT("/sprints/:sprintId", sprintHandler.UpdateSprint)
	api.DELETE("/sprints/:sprintId", sprintHandler.DeleteSprint)
	api.POST("/sprints/:sprintId/userstories", userStoryHandler.AssignUserStoryToSprint)

	// Task Board routes
	api.GET("/sprints/:sprintId/taskboard", taskBoardHandler.GetTaskBoard)

	// Evaluation routes
	api.POST("/evaluations", evaluationHandler.CreateEvaluation)
	api.GET("/evaluations/:id", evaluationHandler.GetEvaluationByID)
	api.GET("/students/:studentId/evaluations", evaluationHandler.GetEvaluationsByStudentID)

	// Metric routes
	api.GET("/sprints/:sprintId/metrics/burndown", metricHandler.GetBurndownChart)
	api.GET("/projects/:projectId/metrics/velocity", metricHandler.GetTeamVelocity)
	api.GET("/sprints/:sprintId/metrics/work-distribution", metricHandler.GetWorkDistribution)

	// Reporting routes
	api.POST("/projects/:projectId/reports/generate", reportingHandler.GenerateProjectReport)

	// --- Admin-Only Routes ---
	admin := e.Group("/api/admin")
	admin.Use(middleware.JWTAuthMiddleware(jwtSecret))
	admin.Use(middleware.AdminAuthMiddleware)

	// Admin user management
	admin.POST("/users", adminHandler.CreateUser)
	admin.POST("/users/admin", adminHandler.CreateAdminUser)
	admin.GET("/users", adminHandler.GetAllUsers)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)

	// Admin project management
	admin.POST("/projects/:id/members", projectHandler.AddMemberToProject)
}
