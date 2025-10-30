package tests

import (
	"strconv"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/routes"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/storage"
	"github.com/buga/API_wrkf/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestRouter holds all the necessary components for running integration tests.
type TestRouter struct {
	Echo           *echo.Echo
	DB             *gorm.DB
	UserService    *services.UserService
	ProjectService *services.ProjectService
	SprintService  *services.SprintService
	UserStoryService *services.UserStoryService
	TaskService    *services.TaskService
}

// NewTestRouter initializes a full application stack in-memory for testing.
func NewTestRouter(cfg *config.AppConfig) (*TestRouter, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := storage.Migrate(db); err != nil {
		return nil, err
	}

	hub := websocket.NewHub()
	go hub.Run()

	// Repositories
	userRepo := storage.NewUserRepository(db)
	projectRepo := storage.NewProjectRepository(db)
	sprintRepo := storage.NewSprintRepository(db)
	taskRepo := storage.NewTaskRepository(db)
	userStoryRepo := storage.NewUserStoryRepository(db)
	notificationRepo := storage.NewNotificationRepository(db)
	rubricRepo := storage.NewRubricRepository(db)
	evaluationRepo := storage.NewEvaluationRepository(db)
	burndownRepo := storage.NewBurndownRepository(db)

	// Services
	notificationService := services.NewNotificationService(notificationRepo)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	projectService := services.NewProjectService(projectRepo, userRepo, userStoryRepo, sprintRepo, taskRepo, notificationService)
	sprintService := services.NewSprintService(sprintRepo)
	taskService := services.NewTaskService(taskRepo, projectService, notificationService, hub)
	userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService)
	rubricService := services.NewRubricService(rubricRepo)
	evaluationService := services.NewEvaluationService(evaluationRepo, rubricRepo)
	burndownService := services.NewBurndownService(burndownRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	sprintHandler := handlers.NewSprintHandler(sprintService)
	taskHandler := handlers.NewTaskHandler(taskService)
	userStoryHandler := handlers.NewUserStoryHandler(userStoryService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	rubricHandler := handlers.NewRubricHandler(rubricService)
	evaluationHandler := handlers.NewEvaluationHandler(evaluationService)
	burndownHandler := handlers.NewBurndownHandler(burndownService)
	websocketHandler := handlers.NewWebsocketHandler(hub, cfg.JWTSecret)

	e := echo.New()
	routes.SetupRoutes(e, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, notificationHandler, rubricHandler, evaluationHandler, burndownHandler, websocketHandler, cfg.JWTSecret)

	return &TestRouter{
		Echo:           e,
		DB:             db,
		UserService:    userService,
		ProjectService: projectService,
		SprintService:  sprintService,
		UserStoryService: userStoryService,
		TaskService:    taskService,
	}, nil
}

// UintToString is a simple helper to convert uint to string for test routes.
func UintToString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
