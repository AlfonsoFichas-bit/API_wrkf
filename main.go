
package main

import (
	"API_wrkf/config"
	"API_wrkf/handlers"
	"API_wrkf/routes"
	"API_wrkf/services"
	"API_wrkf/storage"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load application configuration
	cfg := config.LoadConfig()

	// Create database connection
	db, err := storage.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	// Run database migrations
	if err := storage.Migrate(db); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	fmt.Println("Database migration completed successfully!")

	// --- Initialize Layers with Dependencies ---

	// User components
	userRepo := storage.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)

	// Project components (NEW)
	projectRepo := storage.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	// --- Initialize Echo and Set Up Routes ---
	e := echo.New()
	// Pass all handlers to the router setup function
	routes.SetupRoutes(e, userHandler, projectHandler, cfg.JWTSecret)

	// --- Start Server ---
	fmt.Println("Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
