package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/tests"
	"github.com/stretchr/testify/assert"
)

func TestBurndownChartEndpoint(t *testing.T) {
	// Setup
	cfg := &config.AppConfig{JWTSecret: "test-secret"}
	testRouter, err := tests.NewTestRouter(cfg)
	assert.NoError(t, err)

	// --- 1. Create User and Login to get Token ---
	user := &models.User{
		Nombre:     "Test",
		Correo:     "test@example.com",
		Contraseña: "password",
	}
	err = testRouter.UserService.CreateUser(user)
	assert.NoError(t, err)
	createdUser, err := testRouter.UserService.GetUserByEmail("test@example.com")
	assert.NoError(t, err)

	token, err := testRouter.UserService.Login("test@example.com", "password")
	assert.NoError(t, err)

	// --- 2. Create Project and Add User as Member ---
	project := &models.Project{Name: "Burndown Project"}
	err = testRouter.ProjectService.CreateProject(project, createdUser.ID)
	assert.NoError(t, err)

	_, err = testRouter.ProjectService.AddMemberToProject(project.ID, createdUser.ID, string(models.RoleProductOwner))
	assert.NoError(t, err)

	// --- 3. Create Sprint ---
	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, 4) // 5 day sprint
	sprint := &models.Sprint{
		Name:      "Burndown Sprint",
		StartDate: &startDate,
		EndDate:   &endDate,
	}
	err = testRouter.SprintService.CreateSprint(sprint, project.ID, createdUser.ID)
	assert.NoError(t, err)

	// --- 4. Create User Story and Tasks ---
	userStory := &models.UserStory{Title: "Test US"}
	err = testRouter.UserStoryService.CreateUserStory(userStory, project.ID, createdUser.ID)
	assert.NoError(t, err)

	_, err = testRouter.UserStoryService.AssignUserStoryToSprint(sprint.ID, userStory.ID, createdUser.ID, string(createdUser.Role))
	assert.NoError(t, err)

	task1, err := testRouter.TaskService.CreateTask(&models.Task{Title: "Task 1", StoryPoints: 5}, userStory.ID, createdUser.ID)
	assert.NoError(t, err)
	_, err = testRouter.TaskService.CreateTask(&models.Task{Title: "Task 2", StoryPoints: 3}, userStory.ID, createdUser.ID)
	assert.NoError(t, err)

	// --- 5. Simulate work being done ---
	history := models.TaskHistory{
		TaskID:      task1.ID,
		ChangedByID: createdUser.ID,
		FieldName:   "status",
		NewValue:    "done",
		ChangedAt:   startDate.AddDate(0, 0, 1),
	}
	testRouter.DB.Create(&history)

	// --- 6. Execute Request ---
	req := httptest.NewRequest(http.MethodGet, "/api/sprints/"+tests.UintToString(sprint.ID)+"/burndown", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	testRouter.Echo.ServeHTTP(rec, req)

	// --- 7. Assertions ---
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK")

	var chart services.BurndownChart
	err = json.Unmarshal(rec.Body.Bytes(), &chart)
	assert.NoError(t, err)
	assert.Equal(t, "Burndown Sprint", chart.SprintName)
	assert.Len(t, chart.DataPoints, 5, "Should have 5 data points for a 5-day sprint")

	// Total points = 8. Ideal burn = 2 points/day.
	assert.InDelta(t, 8, chart.DataPoints[0].ActualPoints, 0.01, "Day 1 Actual should be 8")
	assert.InDelta(t, 8, chart.DataPoints[0].IdealPoints, 0.01, "Day 1 Ideal should be 8")

	assert.InDelta(t, 8, chart.DataPoints[1].ActualPoints, 0.01, "Day 2 Actual should be 8")
	assert.InDelta(t, 6, chart.DataPoints[1].IdealPoints, 0.01, "Day 2 Ideal should be 6")

	assert.InDelta(t, 3, chart.DataPoints[2].ActualPoints, 0.01, "Day 3 Actual should be 3")
	assert.InDelta(t, 4, chart.DataPoints[2].IdealPoints, 0.01, "Day 3 Ideal should be 4")
}
