package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/buga/API_wrkf/config"
	_ "github.com/buga/API_wrkf/docs"
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/routes"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/storage"
	"github.com/buga/API_wrkf/websocket"

	"github.com/labstack/echo/v4"
	"github.com/rs/cors"
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

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	cfg := config.LoadConfig()

	db, err := storage.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	if err := storage.Migrate(db); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	// --- Repositories ---
	userRepo := storage.NewUserRepository(db)
	projectRepo := storage.NewProjectRepository(db)
	sprintRepo := storage.NewSprintRepository(db)
	taskRepo := storage.NewTaskRepository(db)
	userStoryRepo := storage.NewUserStoryRepository(db)
	notificationRepo := storage.NewNotificationRepository(db)
	rubricRepo := storage.NewRubricRepository(db)
	burndownRepo := storage.NewBurndownRepository(db)

	// --- Services ---
	notificationService := services.NewNotificationService(notificationRepo)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	projectService := services.NewProjectService(projectRepo, userRepo, userStoryRepo, sprintRepo, taskRepo, notificationService)
	sprintService := services.NewSprintService(sprintRepo)
	taskService := services.NewTaskService(taskRepo, projectService, notificationService, hub)
	userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService)
	rubricService := services.NewRubricService(rubricRepo)
	burndownService := services.NewBurndownService(burndownRepo)

	// --- Handlers ---
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	sprintHandler := handlers.NewSprintHandler(sprintService)
	taskHandler := handlers.NewTaskHandler(taskService)
	userStoryHandler := handlers.NewUserStoryHandler(userStoryService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	rubricHandler := handlers.NewRubricHandler(rubricService)
	burndownHandler := handlers.NewBurndownHandler(burndownService)
	websocketHandler := handlers.NewWebsocketHandler(hub, cfg.JWTSecret)

	// --- Admin User ---
	createAdminUserIfNeeded(userService, cfg.Admin)

	// --- Echo Setup ---
	e := echo.New()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	})
	e.Use(echo.WrapMiddleware(c.Handler))

	routes.SetupRoutes(e, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, notificationHandler, rubricHandler, burndownHandler, websocketHandler, cfg.JWTSecret)

	// --- Start Server ---
	fmt.Println("Iniciando el servidor en el puerto 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}

func createAdminUserIfNeeded(userService *services.UserService, adminCfg *config.AdminConfig) {
	_, err := userService.GetUserByEmail(adminCfg.Email)
	if err == nil {
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error al checkear admin user: %v", err)
		return
	}
	admin := &models.User{
		Nombre:          adminCfg.Nombre,
		ApellidoPaterno: "Admin",
		ApellidoMaterno: "User",
		Correo:          adminCfg.Email,
		Contraseña:      adminCfg.Password,
	}
	if err := userService.CreateAdminUser(admin); err != nil {
		log.Fatalf("No se pudo crear el usuario administrador.: %v", err)
	}
	fmt.Println("Usuario administrador creado correctamente")
}
