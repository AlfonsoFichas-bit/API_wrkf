package tests

import (
	"log"
	"testing"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/routes"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/storage"
	"github.com/buga/API_wrkf/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestApp holds all the components needed to run an integration test.
type TestApp struct {
	DB                  *gorm.DB
	Router              *echo.Echo
	UserService         *services.UserService
	ProjectService      *services.ProjectService
	SprintService       *services.SprintService
	UserStoryService    *services.UserStoryService
	TaskService         *services.TaskService
	NotificationService *services.NotificationService
	RubricService       services.RubricService
	EvaluationService   *services.EvaluationService
	EventService        *services.EventService
	ExportService       *services.ExportService
}

// SetupTestApp initializes a full application stack for integration testing.
func SetupTestApp() *TestApp {
	cfg := config.LoadConfig()

	db, err := storage.NewTestConnection()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	// Initialize components
	userRepo := storage.NewUserRepository(db)
	projectRepo := storage.NewProjectRepository(db)
	userStoryRepo := storage.NewUserStoryRepository(db)
	sprintRepo := storage.NewSprintRepository(db)
	taskRepo := storage.NewTaskRepository(db)
	notificationRepo := storage.NewNotificationRepository(db)
	rubricRepo := storage.NewRubricRepository(db)
	reportingRepo := storage.NewReportingRepository(db)
	evalRepo := storage.NewEvaluationRepository(db)
	eventRepo := storage.NewEventRepository(db)

	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	notificationService := services.NewNotificationService(notificationRepo)
	projectService := services.NewProjectService(projectRepo, userRepo, userStoryRepo, sprintRepo, taskRepo, notificationService)
	sprintService := services.NewSprintService(sprintRepo)
	taskService := services.NewTaskService(taskRepo, projectService, notificationService)
	userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService)
	rubricService := services.NewRubricService(rubricRepo)
	reportingService := services.NewReportingService(reportingRepo, userStoryRepo, sprintRepo)
	evaluationService := services.NewEvaluationService(evalRepo, taskRepo, rubricRepo, projectService)
	eventService := services.NewEventService(eventRepo, projectService)
	exportService := services.NewExportService(projectRepo, userStoryRepo, taskRepo) // <-- NEW

	// WebSocket Manager
	wsManager := websocket.NewWebSocketManager()
	go wsManager.Run()

	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	sprintHandler := handlers.NewSprintHandler(sprintService)
	userStoryHandler := handlers.NewUserStoryHandler(userStoryService)
	taskHandler := handlers.NewTaskHandler(taskService, wsManager, userService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	rubricHandler := handlers.NewRubricHandler(rubricService)
	reportingHandler := handlers.NewReportingHandler(reportingService)
	evaluationHandler := handlers.NewEvaluationHandler(evaluationService)
	eventHandler := handlers.NewEventHandler(eventService)
	exportHandler := handlers.NewExportHandler(exportService) // <-- NEW
	webSocketHandler := websocket.NewWebSocketHandler(wsManager, cfg.JWTSecret, userService, projectService)

	router := echo.New()
	routes.SetupRoutes(router, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, notificationHandler, rubricHandler, reportingHandler, evaluationHandler, eventHandler, exportHandler, cfg.JWTSecret)
	router.GET("/ws", webSocketHandler.HandleConnection)

	return &TestApp{
		DB:                  db,
		Router:              router,
		UserService:         userService,
		ProjectService:      projectService,
		SprintService:       sprintService,
		UserStoryService:    userStoryService,
		TaskService:         taskService,
		NotificationService: notificationService,
		RubricService:       rubricService,
		EvaluationService:   evaluationService,
		EventService:        eventService,
		ExportService:       exportService, // <-- NEW
	}
}

// TeardownTestApp cleans up resources used by the test application.
func TeardownTestApp(app *TestApp) {
	sqlDB, _ := app.DB.DB()
	sqlDB.Close()
}

// Helper functions to create test data easily

func CreateTestUser(t *testing.T, app *TestApp, email, role string) (*models.User, string) {
	user := &models.User{
		Nombre: "Test",
		Correo: email,
		Role:   role,
	}
	err := app.DB.Create(user).Error
	require.NoError(t, err)

	token, err := app.UserService.GenerateJWT(user.ID)
	require.NoError(t, err)
	return user, token
}

func CreateTestProject(t *testing.T, app *TestApp, name string, creatorID uint) *models.Project {
	project := &models.Project{
		Name:        name,
		CreatedByID: creatorID,
	}
	err := app.DB.Create(project).Error
	require.NoError(t, err)
	return project
}

func AddUserToProject(t *testing.T, app *TestApp, projectID, userID uint, role string) {
	member := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	err := app.DB.Create(member).Error
	require.NoError(t, err)
}

func CreateTestUserStory(t *testing.T, app *TestApp, title string, projectID uint) *models.UserStory {
	us := &models.UserStory{
		Title:     title,
		ProjectID: projectID,
	}
	err := app.DB.Create(us).Error
	require.NoError(t, err)
	return us
}

func CreateTestTask(t *testing.T, app *TestApp, title string, userStoryID uint, assignedToID uint) *models.Task {
	task := &models.Task{
		Title:        title,
		UserStoryID:  userStoryID,
		AssignedToID: &assignedToID,
	}
	err := app.DB.Create(task).Error
	require.NoError(t, err)
	return task
}

func CreateTestRubric(t *testing.T, app *TestApp, projectID, creatorID uint, name string) *models.Rubric {
	rubric := &models.Rubric{
		Name:        name,
		ProjectID:   projectID,
		CreatedByID: creatorID,
		Status:      models.RubricStatusActive,
		Criteria: []models.RubricCriterion{
			{Title: "Criterion 1", MaxPoints: 5},
			{Title: "Criterion 2", MaxPoints: 5},
		},
	}
	err := app.DB.Create(rubric).Error
	require.NoError(t, err)
	return rubric
}
