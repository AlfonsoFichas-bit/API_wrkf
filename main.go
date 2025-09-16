package main

import (
	"API_wrkf/config"
	"API_wrkf/handlers"
	"API_wrkf/models"
	"API_wrkf/routes"
	"API_wrkf/services"
	"API_wrkf/storage"
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// createAdminUserIfNeeded checks if an admin user exists and creates one if it doesn't.
func createAdminUserIfNeeded(userService *services.UserService, adminCfg *config.AdminConfig) {
	// Check if a user with the admin email already exists
	_, err := userService.GetUserByEmail(adminCfg.Email)
	if err == nil {
		// Admin user already exists, no action needed
		fmt.Println("Admin user already exists.")
		return
	}

	// If the error is anything other than "record not found", log it and stop
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking for admin user: %v", err)
		return
	}

	// Admin user does not exist, so create one
	fmt.Println("Admin user not found, creating a new one...")
	admin := &models.User{
		Nombre:          adminCfg.Nombre,
		ApellidoPaterno: "Admin", // Default value
		ApellidoMaterno: "User",  // Default value
		Correo:          adminCfg.Email,
		Contraseña:      adminCfg.Password,
	}

	if err := userService.CreateAdminUser(admin); err != nil {
		log.Fatalf("Could not create admin user: %v", err)
	}

	fmt.Println("Admin user created successfully!")
}

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

	// Create admin user if it doesn't exist
	createAdminUserIfNeeded(userService, cfg.Admin)

	// Project components (NEW)
	projectRepo := storage.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	// Sprint components
	sprintRepo := storage.NewSprintRepository(db)
	sprintService := services.NewSprintService(sprintRepo)
	sprintHandler := handlers.NewSprintHandler(sprintService)

	// --- Initialize Echo and Set Up Routes ---
	e := echo.New()
	// Pass all handlers to the router setup function
	routes.SetupRoutes(e, userHandler, projectHandler, sprintHandler, cfg.JWTSecret)

	// --- Start Server ---
	fmt.Println("Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
