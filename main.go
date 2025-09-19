package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/buga/API_wrkf/config"
	_ "github.com/buga/API_wrkf/docs" // This line is needed for swag to find your docs
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/routes"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/storage"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// @title           API-wrkf Project Management API
// @version         1.0
// @description     This is a comprehensive API for a Scrum-based project management platform.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
// @description "Type 'Bearer' followed by a space and the JWT token. Example: 'Bearer {token}"

// createAdminUserIfNeeded comprueba si existe un usuario administrador si no crea uno.
func createAdminUserIfNeeded(userService *services.UserService, adminCfg *config.AdminConfig) {
	_, err := userService.GetUserByEmail(adminCfg.Email)
	if err == nil {
		fmt.Println("Usuario administrador ya existe.")
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking for admin user: %v", err)
		return
	}

	fmt.Println("Admin user not found, creating a new one...")
	admin := &models.User{
		Nombre:          adminCfg.Nombre,
		ApellidoPaterno: "Admin",
		ApellidoMaterno: "User",
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

	// Sprint, Task, and UserStory Repositories
	sprintRepo := storage.NewSprintRepository(db)
	taskRepo := storage.NewTaskRepository(db)
	userStoryRepo := storage.NewUserStoryRepository(db)

	// Project Service (now with more dependencies)
	projectRepo := storage.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo, userStoryRepo, sprintRepo, taskRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	// Other Services
	sprintService := services.NewSprintService(sprintRepo)
	taskService := services.NewTaskService(taskRepo, projectService)
	userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService)

	// Handlers
	sprintHandler := handlers.NewSprintHandler(sprintService)
	taskHandler := handlers.NewTaskHandler(taskService)
	userStoryHandler := handlers.NewUserStoryHandler(userStoryService)

	// --- Initialize Echo and Set Up Routes ---
	e := echo.New()
	routes.SetupRoutes(e, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, cfg.JWTSecret)

	// --- Start Server ---
	fmt.Println("Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
