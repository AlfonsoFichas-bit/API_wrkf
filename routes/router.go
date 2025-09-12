package routes

import (
	"API_wrkf/handlers"

	"API_wrkf/middleware"

	"github.com/labstack/echo/v4"
)

// SetupRoutes configures the application routes.
// It now accepts the jwtSecret to create the auth middleware.
func SetupRoutes(e *echo.Echo, userHandler *handlers.UserHandler, jwtSecret string) {
	// --- Public Routes ---
	e.POST("/login", userHandler.Login)
	e.POST("/users", userHandler.CreateUser)

	// --- Protected Routes ---
	// Create a group for routes that require JWT authentication.
	r := e.Group("")

	// Create the JWT middleware using the factory and the provided secret.
	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)
	r.Use(authMiddleware)

	// All routes defined in this group are now protected.
	r.GET("/users/:id", userHandler.GetUser)
}
