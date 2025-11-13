package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buga/API_wrkf/handlers"
	"github.com/buga/API_wrkf/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskAssignmentAndReassignment(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// 1. Setup: Create users, project, story, and a task
	admin, adminToken := CreateTestUser(t, testApp, "task_admin@test.com", "admin")
	userA, _ := CreateTestUser(t, testApp, "user_a@test.com", "developer")
	userB, _ := CreateTestUser(t, testApp, "user_b@test.com", "developer")

	project := CreateTestProject(t, testApp, "Assignment Project", admin.ID)
	AddUserToProject(t, testApp, project.ID, userA.ID, "developer")
	AddUserToProject(t, testApp, project.ID, userB.ID, "developer")

	us := CreateTestUserStory(t, testApp, "Assignment Story", project.ID)
	task := &models.Task{
		Title:       "Test Task",
		UserStoryID: us.ID,
		Status:      models.StatusTodo,
		CreatedByID: admin.ID,
	}
	err := testApp.DB.Create(task).Error
	require.NoError(t, err)

	// 2. Assign to User A
	assignReq := handlers.AssignTaskRequest{UserID: userA.ID}
	body, _ := json.Marshal(assignReq)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/tasks/%d/assign", task.ID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+adminToken)
	rec := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// Verify assignment to User A
	taskAfterAssign, err := testApp.TaskService.GetTaskByID(task.ID)
	require.NoError(t, err)
	require.NotNil(t, taskAfterAssign.AssignedToID)
	assert.Equal(t, userA.ID, *taskAfterAssign.AssignedToID)

	// 3. Reassign to User B
	reassignReq := handlers.AssignTaskRequest{UserID: userB.ID}
	body, _ = json.Marshal(reassignReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/tasks/%d/assign", task.ID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+adminToken)
	rec = httptest.NewRecorder()
	testApp.Router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// 4. Verify reassignment to User B
	taskAfterReassign, err := testApp.TaskService.GetTaskByID(task.ID)
	require.NoError(t, err)
	require.NotNil(t, taskAfterReassign.AssignedToID)
	assert.Equal(t, userB.ID, *taskAfterReassign.AssignedToID, "Task should be reassigned to User B")
}

func TestKanbanUpdateKeepsAssignment(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// 1. Setup
	user, userToken := CreateTestUser(t, testApp, "kanban_user@test.com", "developer")
	project := CreateTestProject(t, testApp, "Kanban Project", user.ID)
	us := CreateTestUserStory(t, testApp, "Kanban Story", project.ID)
	task := &models.Task{
		Title:        "Kanban Task",
		UserStoryID:  us.ID,
		Status:       models.StatusTodo,
		AssignedToID: &user.ID,
		CreatedByID:  user.ID,
	}
	err := testApp.DB.Create(task).Error
	require.NoError(t, err)

	// 2. Action: Update status via API
	updateData := map[string]string{"status": string(models.StatusInProgress)}
	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/tasks/%d/status", task.ID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+userToken)
	rec := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// 3. Verification
	updatedTask, err := testApp.TaskService.GetTaskByID(task.ID)
	require.NoError(t, err)

	assert.Equal(t, models.StatusInProgress, updatedTask.Status, "Status should have been updated.")
	require.NotNil(t, updatedTask.AssignedToID, "AssignedToID should not be nil after status update.")
	assert.Equal(t, user.ID, *updatedTask.AssignedToID, "Task should remain assigned to the user.")
}

func TestTaskCreationWithAssignment(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// 1. Setup
	admin, adminToken := CreateTestUser(t, testApp, "task_admin@test.com", "admin")
	user, _ := CreateTestUser(t, testApp, "assignee@test.com", "developer")
	project := CreateTestProject(t, testApp, "Creation Project", admin.ID)
	AddUserToProject(t, testApp, project.ID, user.ID, "developer")
	us := CreateTestUserStory(t, testApp, "Creation Story", project.ID)

	// 2. Action: Create a new task with a pre-assigned user
	taskData := models.Task{
		Title:        "Pre-assigned Task",
		UserStoryID:  us.ID,
		AssignedToID: &user.ID,
	}
	body, _ := json.Marshal(taskData)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/userstories/%d/tasks", us.ID), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+adminToken)
	rec := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rec, req)

	// Assertions
	require.Equal(t, http.StatusCreated, rec.Code)
	var createdTask models.Task
	err := json.Unmarshal(rec.Body.Bytes(), &createdTask)
	require.NoError(t, err)

	// Verification
	fetchedTask, err := testApp.TaskService.GetTaskByID(createdTask.ID)
	require.NoError(t, err)
	require.NotNil(t, fetchedTask.AssignedToID, "Task should be assigned upon creation")
	assert.Equal(t, user.ID, *fetchedTask.AssignedToID, "Task should be assigned to the correct user")
}

func TestDeleteTaskWithHistory(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// 1. Setup
	user, userToken := CreateTestUser(t, testApp, "delete_user@test.com", "developer")
	project := CreateTestProject(t, testApp, "Delete Project", user.ID)
	us := CreateTestUserStory(t, testApp, "Delete Story", project.ID)
	task := &models.Task{
		Title:       "Delete Task",
		UserStoryID: us.ID,
		Status:      models.StatusTodo,
		CreatedByID: user.ID,
	}
	err := testApp.DB.Create(task).Error
	require.NoError(t, err)

	// Create a history record for the task
	history := &models.TaskHistory{
		TaskID:      task.ID,
		ChangedByID: user.ID,
		FieldName:   "title",
		OldValue:    "Old",
		NewValue:    "New",
	}
	err = testApp.DB.Create(history).Error
	require.NoError(t, err)

	// 2. Action: Delete the task via API
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/tasks/%d", task.ID), nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+userToken)
	rec := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rec, req)

	// 3. Verification
	assert.Equal(t, http.StatusNoContent, rec.Code, "Deleting the task should be successful")

	// Verify the task is actually gone
	_, err = testApp.TaskService.GetTaskByID(task.ID)
	assert.Error(t, err, "Task should not be found after deletion")
}
