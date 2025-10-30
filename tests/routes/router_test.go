package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/storage"
	"github.com/buga/API_wrkf/tests"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginEndpoint(t *testing.T) {
	cfg := &config.AppConfig{JWTSecret: "test-secret-login"}
	testRouter, err := tests.NewTestRouter(cfg)
	assert.NoError(t, err)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Nombre:     "testuser",
		Correo:     "test@example.com",
		Contraseña: string(hashedPassword),
	}
	testRouter.DB.Create(&user)

	t.Run("Successful Login", func(t *testing.T) {
		loginCredentials := map[string]string{
			"email":    "test@example.com",
			"password": password,
		}
		jsonBody, _ := json.Marshal(loginCredentials)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		testRouter.Echo.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to unmarshal response body")
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])
	})
}

func TestProjectBoardEndpoint(t *testing.T) {
	cfg := &config.AppConfig{JWTSecret: "test-secret-board"}
	testRouter, err := tests.NewTestRouter(cfg)
	assert.NoError(t, err)

	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Nombre: "testuser", Correo: "boarduser@example.com", Contraseña: string(hashedPassword)}
	testRouter.DB.Create(&user)

	project := models.Project{Name: "Board Project", CreatedByID: user.ID}
	testRouter.DB.Create(&project)
	userStory := models.UserStory{Title: "Board Story", ProjectID: project.ID, CreatedByID: user.ID}
	testRouter.DB.Create(&userStory)
	tasks := []models.Task{
		{Title: "Todo Task", UserStoryID: userStory.ID, Status: string(models.StatusTodo), CreatedByID: user.ID},
		{Title: "In Progress Task", UserStoryID: userStory.ID, Status: string(models.StatusInProgress), CreatedByID: user.ID},
	}
	testRouter.DB.Create(&tasks)

	token, err := testRouter.UserService.Login("boarduser@example.com", password)
	assert.NoError(t, err)

	t.Run("Get Project Board", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/projects/%d/board", project.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		// Manually build context and handler for this specific case
		c := testRouter.Echo.NewContext(req, rec)
		c.Set("userID", float64(user.ID))
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", project.ID))

		userRepo := storage.NewUserRepository(testRouter.DB)
		userStoryRepo := storage.NewUserStoryRepository(testRouter.DB)
		sprintRepo := storage.NewSprintRepository(testRouter.DB)
		taskRepo := storage.NewTaskRepository(testRouter.DB)
		notificationRepo := storage.NewNotificationRepository(testRouter.DB)
		notificationService := services.NewNotificationService(notificationRepo)
		projectRepo := storage.NewProjectRepository(testRouter.DB)
		projectService := services.NewProjectService(projectRepo, userRepo, userStoryRepo, sprintRepo, taskRepo, notificationService)
		handler := handlers.NewProjectHandler(projectService)

		if assert.NoError(t, handler.GetProjectBoard(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var response map[string][]models.Task
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err, "Failed to unmarshal board response")
			assert.Len(t, response["todo"], 1)
			assert.Len(t, response["in_progress"], 1)
		}
	})
}
