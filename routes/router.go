
package routes

import (
	"API_wrkf/handlers"
	"API_wrkf/middleware"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures the application routes.
func SetupRoutes(e *echo.Echo, userHandler *handlers.UserHandler, projectHandler *handlers.ProjectHandler, jwtSecret string) {
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
	api.PUT("/projects/:id", projectHandler.UpdateProject)       // <-- NEW
	api.DELETE("/projects/:id", projectHandler.DeleteProject) // <-- NEW

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
