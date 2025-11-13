package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKanbanStatusUpdateRegression(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// 1. Setup: Create a user, project, story, and an assigned task
	user, userToken := CreateTestUser(t, testApp, "kanban_user@test.com", "developer")
	project := CreateTestProject(t, testApp, "Kanban Project", user.ID)
	us := CreateTestUserStory(t, testApp, "Kanban Story", project.ID)

	task := &models.Task{
		Title:        "Kanban Task",
		UserStoryID:  us.ID,
		Status:       models.StatusTodo,
		AssignedToID: &user.ID, // Assign the task
		CreatedByID:  user.ID,
	}
	err := testApp.DB.Create(task).Error
	require.NoError(t, err)

	// Verify initial state
	initialTask, err := testApp.TaskService.GetTaskByID(task.ID)
	require.NoError(t, err)
	require.NotNil(t, initialTask.AssignedToID, "Task should be assigned initially")
	assert.Equal(t, user.ID, *initialTask.AssignedToID)

	// 2. Action: Call the endpoint to update the task's status
	updateData := map[string]string{"status": string(models.StatusInProgress)}
	body, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/tasks/%d/status", task.ID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+userToken)
	rec := httptest.NewRecorder()

	testApp.Router.ServeHTTP(rec, req)

	// Assert API call was successful
	require.Equal(t, http.StatusOK, rec.Code, "API call to update status should succeed")

	// 3. Verification: Fetch the task again and check its fields
	updatedTask, err := testApp.TaskService.GetTaskByID(task.ID)
	require.NoError(t, err, "Fetching the task after update should not fail")

	// This is the critical part of the test
	assert.Equal(t, models.StatusInProgress, updatedTask.Status, "Task status should be updated")
	assert.NotNil(t, updatedTask.AssignedToID, "AssignedToID should NOT be nil after status update")
	assert.Equal(t, user.ID, *updatedTask.AssignedToID, "Task should remain assigned to the same user")
}
