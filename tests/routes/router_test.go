package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/routes"
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

    userStoryService := services.NewUserStoryService(userStoryRepo, projectService, sprintService)
    userStoryHandler := handlers.NewUserStoryHandler(userStoryService)

    taskService := services.NewTaskService(taskRepo, projectService, notificationService)
    taskHandler := handlers.NewTaskHandler(taskService)

    notificationHandler := handlers.NewNotificationHandler(notificationService)

    rubricRepo := storage.NewRubricRepository(db)
    rubricService := services.NewRubricService(rubricRepo)
    rubricHandler := handlers.NewRubricHandler(rubricService)


	e := echo.New()
	evaluationRepo := storage.NewEvaluationRepository(db)
	evaluationService := services.NewEvaluationService(evaluationRepo, taskService, projectService)
	evaluationHandler := handlers.NewEvaluationHandler(evaluationService)

	routes.SetupRoutes(e, userHandler, projectHandler, sprintHandler, userStoryHandler, taskHandler, notificationHandler, rubricHandler, evaluationHandler, cfg.JWTSecret)

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

func TestProjectBoardEndpoint(t *testing.T) {
	e, db := setupTestApp()

	// --- Setup Test Data ---
	// User
	user := models.User{Nombre: "testuser", Correo: "boarduser@example.com", Contraseña: "password"}
	db.Create(&user)
	// Project
	project := models.Project{Name: "Board Project", CreatedByID: user.ID}
	db.Create(&project)
	// User Story
	userStory := models.UserStory{Title: "Board Story", ProjectID: project.ID, CreatedByID: user.ID}
	db.Create(&userStory)
	// Tasks
	task1 := models.Task{Title: "Todo Task", UserStoryID: userStory.ID, Status: string(models.StatusTodo), CreatedByID: user.ID}
	task2 := models.Task{Title: "In Progress Task", UserStoryID: userStory.ID, Status: string(models.StatusInProgress), CreatedByID: user.ID}
	db.Create(&task1)
	db.Create(&task2)

	// --- Test Case ---
	t.Run("Get Project Board", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/projects/"+strconv.Itoa(int(project.ID))+"/board", nil)
		rec := httptest.NewRecorder()

		// Create a fake context and JWT for the request
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(project.ID)))
		c.Set("userID", float64(user.ID))

		// Manually call the handler because middleware is tricky to run in tests
		handler := handlers.NewProjectHandler(services.NewProjectService(
			storage.NewProjectRepository(db),
			storage.NewUserRepository(db),
			storage.NewUserStoryRepository(db),
			storage.NewSprintRepository(db),
			storage.NewTaskRepository(db),
			services.NewNotificationService(storage.NewNotificationRepository(db)),
		))

		if assert.NoError(t, handler.GetProjectBoard(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response map[string][]models.Task
			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Contains(t, response, "todo")
			assert.Contains(t, response, "in_progress")
			assert.Contains(t, response, "in_review")
			assert.Contains(t, response, "done")
			assert.Len(t, response["todo"], 1)
			assert.Len(t, response["in_progress"], 1)
			assert.Len(t, response["done"], 0)
			assert.Equal(t, "Todo Task", response["todo"][0].Title)
			assert.Equal(t, "In Progress Task", response["in_progress"][0].Title)
		}
	})
}

func TestCreateEvaluationEndpoint(t *testing.T) {
	e, db := setupTestApp()

	// --- Setup Test Data ---
	user := models.User{Nombre: "evaluator", Correo: "evaluator@example.com", Contraseña: "password"}
	db.Create(&user)
	project := models.Project{Name: "Evaluation Project", CreatedByID: user.ID}
	db.Create(&project)
	// Assign user to project with a valid role for evaluation
	projectMember := models.ProjectMember{ProjectID: project.ID, UserID: user.ID, Role: string(models.RoleProductOwner)}
	db.Create(&projectMember)
	userStory := models.UserStory{Title: "Evaluation Story", ProjectID: project.ID, CreatedByID: user.ID}
	db.Create(&userStory)
	task := models.Task{Title: "Deliverable Task", UserStoryID: userStory.ID, IsDeliverable: true, CreatedByID: user.ID}
	db.Create(&task)

	// --- Test Case ---
	t.Run("Create Evaluation for Deliverable", func(t *testing.T) {
		evaluationData := map[string]interface{}{
			"evaluatorId": user.ID,
			"score":       95.5,
			"comments":    "Excellent work!",
		}
		jsonBody, _ := json.Marshal(evaluationData)

		req := httptest.NewRequest(http.MethodPost, "/api/tasks/"+strconv.Itoa(int(task.ID))+"/evaluations", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Create a fake context
		c := e.NewContext(req, rec)
		c.Set("userID", float64(user.ID))
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(task.ID)))

		// Instantiate the handler with real dependencies
		userRepo := storage.NewUserRepository(db)
		projectRepo := storage.NewProjectRepository(db)
		taskRepo := storage.NewTaskRepository(db)
		notificationRepo := storage.NewNotificationRepository(db)
		notificationService := services.NewNotificationService(notificationRepo)
		projectService := services.NewProjectService(projectRepo, userRepo, nil, nil, nil, notificationService) // Some dependencies can be nil if not used in this specific test path
		taskService := services.NewTaskService(taskRepo, projectService, notificationService)
		evalRepo := storage.NewEvaluationRepository(db)
		evalService := services.NewEvaluationService(evalRepo, taskService, projectService)
		evalHandler := handlers.NewEvaluationHandler(evalService)

		// This will fail initially because the handler returns 501 Not Implemented
		if assert.NoError(t, evalHandler.CreateEvaluation(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response models.Evaluation
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, uint(task.ID), response.TaskID)
			assert.Equal(t, "Excellent work!", response.Comments)
		}
	})
}
