package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp() (*echo.Echo, *gorm.DB) {
	// Setup in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Run migrations
	err = storage.Migrate(db)
	if err != nil {
		panic("failed to migrate database")
	}

	// --- Initialize layers ---
	cfg := &config.AppConfig{JWTSecret: "test_secret"}
	userRepo := storage.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)

	// Dependencies for other handlers (can be nil if not used in the test)
	projectRepo := storage.NewProjectRepository(db)
    notificationRepo := storage.NewNotificationRepository(db)
    notificationService := services.NewNotificationService(notificationRepo)
    userStoryRepo := storage.NewUserStoryRepository(db)
    sprintRepo := storage.NewSprintRepository(db)
    taskRepo := storage.NewTaskRepository(db)
	projectService := services.NewProjectService(projectRepo, userRepo, userStoryRepo, sprintRepo, taskRepo, notificationService)
	projectHandler := handlers.NewProjectHandler(projectService)

    sprintService := services.NewSprintService(sprintRepo)
    sprintHandler := handlers.NewSprintHandler(sprintService)

    sprintService = services.NewSprintService(sprintRepo)
    userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService)
    userStoryHandler := handlers.NewUserStoryHandler(userStoryService)

    taskService := services.NewTaskService(taskRepo, projectService, notificationService)
    taskHandler := handlers.NewTaskHandler(taskService)

    notificationHandler := handlers.NewNotificationHandler(notificationService)

    rubricRepo := storage.NewRubricRepository(db)
    rubricService := services.NewRubricService(rubricRepo)
    rubricHandler := handlers.NewRubricHandler(rubricService)


	e := echo.New()
	SetupRoutes(e, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, notificationHandler, rubricHandler, cfg.JWTSecret)

	return e, db
}

func TestLoginEndpoint(t *testing.T) {
	e, db := setupTestApp()

	// --- Setup Test Data ---
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Nombre:     "testuser",
		Correo:     "test@example.com",
		Contraseña: string(hashedPassword),
	}
	db.Create(&user)

	// --- Test Case: Successful Login ---
	t.Run("Successful Login", func(t *testing.T) {
		loginCredentials := map[string]string{
			"correo":     "test@example.com",
			"contraseña": password,
		}
		jsonBody, _ := json.Marshal(loginCredentials)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])
	})
}
