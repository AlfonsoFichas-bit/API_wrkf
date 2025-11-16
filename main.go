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

	// Repositories
	userRepo := storage.NewUserRepository(db)
	notificationRepo := storage.NewNotificationRepository(db)
	sprintRepo := storage.NewSprintRepository(db)
	taskRepo := storage.NewTaskRepository(db)
	userStoryRepo := storage.NewUserStoryRepository(db)
	projectRepo := storage.NewProjectRepository(db)
	rubricRepo := storage.NewRubricRepository(db)
	reportingRepo := storage.NewReportingRepository(db)
	evalRepo := storage.NewEvaluationRepository(db)
	eventRepo := storage.NewEventRepository(db)
	activityRepo := storage.NewActivityRepository(db)
	deadlineRepo := storage.NewDeadlineRepository(db)
	teamRepo := storage.NewTeamRepository(db) // New

	// Services
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	notificationService := services.NewNotificationService(notificationRepo)
	teamService := services.NewTeamService(teamRepo) // New
	projectService := services.NewProjectService(projectRepo, userRepo, userStoryRepo, sprintRepo, taskRepo, notificationService, teamService) // Updated
	reportingService := services.NewReportingService(reportingRepo, userStoryRepo, sprintRepo)
	activityService := services.NewActivityService(activityRepo, projectService)
	sprintService := services.NewSprintService(sprintRepo, reportingService, activityService) // Updated
	taskService := services.NewTaskService(taskRepo, projectService, notificationService, activityService)
	userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService, activityService) // Updated
	rubricService := services.NewRubricService(rubricRepo)
	evaluationService := services.NewEvaluationService(evalRepo, taskRepo, rubricRepo, projectService)
	eventService := services.NewEventService(eventRepo, projectService)
	exportService := services.NewExportService(projectRepo, userStoryRepo, taskRepo)
	deadlineService := services.NewDeadlineService(deadlineRepo, projectService) // New

	// WebSocket Manager (initialized before handlers that need it)
	wsManager := websocket.NewWebSocketManager()
	go wsManager.Run()

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	projectHandler := handlers.NewProjectHandler(projectService)
	sprintHandler := handlers.NewSprintHandler(sprintService)
	userStoryHandler := handlers.NewUserStoryHandler(userStoryService)
	rubricHandler := handlers.NewRubricHandler(rubricService)
	reportingHandler := handlers.NewReportingHandler(reportingService)
	evaluationHandler := handlers.NewEvaluationHandler(evaluationService, projectRepo)
	eventHandler := handlers.NewEventHandler(eventService)
	exportHandler := handlers.NewExportHandler(exportService)
	activityHandler := handlers.NewActivityHandler(activityService)
	deadlineHandler := handlers.NewDeadlineHandler(deadlineService, teamService) // Updated
	taskHandler := handlers.NewTaskHandler(taskService, wsManager, userService)
	webSocketHandler := websocket.NewWebSocketHandler(wsManager, cfg.JWTSecret, userService, projectService)

	// Final setup
	createAdminUserIfNeeded(userService, cfg.Admin)

	// --- Inicializar Echo y configurar routes ---
	e := echo.New()
	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // URL del frontend en desarrollo
			"http://localhost:3000", // Si usas otro puerto
			"http://127.0.0.1:5173",
			"http://0.0.0.0:8000",
			"http://127.0.0.1:8000",
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

	// WebSocket route
	e.GET("/ws", webSocketHandler.HandleConnection)

	routes.SetupRoutes(e, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, notificationHandler, rubricHandler, reportingHandler, evaluationHandler, eventHandler, exportHandler, activityHandler, deadlineHandler, cfg.JWTSecret)

	// --- Iniciar Servidor ---
	fmt.Println("Iniciando el servidor en el puerto 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}
