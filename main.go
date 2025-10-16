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

// createAdminUserIfNeeded comprueba si existe un usuario administrador si no crea uno.
func createAdminUserIfNeeded(userService *services.UserService, adminCfg *config.AdminConfig) {
	_, err := userService.GetUserByEmail(adminCfg.Email)
	if err == nil {
		fmt.Println("Usuario administrador ya existe.")
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error al checkear admin user: %v", err)
		return
	}

	fmt.Println("No se ha encontrado el usuario administrador, creando uno nuevo...")
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

func main() {
	// Cargar configuración de la aplicación
	cfg := config.LoadConfig()

	// Crear conexión a la db
	db, err := storage.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	// Ejecutar migraciones de db
	if err := storage.Migrate(db); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}

	fmt.Println("Database migration completed successfully!")

	// --- Inicializar capas con dependencias ---

	// User components
	userRepo := storage.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)

	// Cargar configuración de la aplicación
	createAdminUserIfNeeded(userService, cfg.Admin)

	// Notification components
	notificationRepo := storage.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Rubric components
	rubricRepo := storage.NewRubricRepository(db)
	rubricService := services.NewRubricService(rubricRepo)
	rubricHandler := handlers.NewRubricHandler(rubricService)

	// Sprint, Task, and UserStory Repositories
	sprintRepo := storage.NewSprintRepository(db)
	taskRepo := storage.NewTaskRepository(db)
	userStoryRepo := storage.NewUserStoryRepository(db)

	// Project Service (now with more dependencies)
	projectRepo := storage.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo, userStoryRepo, sprintRepo, taskRepo, notificationService) // Inyectar notificationService
	projectHandler := handlers.NewProjectHandler(projectService)

	// Other Services
	sprintService := services.NewSprintService(sprintRepo)
	taskService := services.NewTaskService(taskRepo, projectService, notificationService)
	userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService)

	// Handlers
	sprintHandler := handlers.NewSprintHandler(sprintService)
	taskHandler := handlers.NewTaskHandler(taskService)
	userStoryHandler := handlers.NewUserStoryHandler(userStoryService)

	// --- Inicializar Echo y configurar routes ---
	e := echo.New()
	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // URL del frontend en desarrollo
			"http://localhost:3000", // Si usas otro puerto
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},
		AllowCredentials: true,
	})
	e.Use(echo.WrapMiddleware(c.Handler))
	routes.SetupRoutes(e, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, notificationHandler, rubricHandler, cfg.JWTSecret)

	// --- Iniciar Servidor ---
	fmt.Println("Iniciando el servidor en el puerto 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}
