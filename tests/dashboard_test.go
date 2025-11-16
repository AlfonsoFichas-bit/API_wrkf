package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardEndpoints(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// Create users with different roles
	student, studentToken := CreateTestUser(t, testApp, "student_dash@test.com", "user")
	student2, _ := CreateTestUser(t, testApp, "student2_dash@test.com", "user")
	admin, adminToken := CreateTestUser(t, testApp, "admin_dash@test.com", "admin")

	// Create a project and add users
	project := CreateTestProject(t, testApp, "Dashboard Test Project", admin.ID)

	// The creator (admin) is automatically added as a product_owner.
	// We need to find that membership record and update the role to 'docente' for the test.
	var creatorMember models.ProjectMember
	err := testApp.DB.Where("project_id = ? AND user_id = ?", project.ID, admin.ID).First(&creatorMember).Error
	require.NoError(t, err)
	creatorMember.Role = "docente"
	err = testApp.DB.Save(&creatorMember).Error
	require.NoError(t, err)

	AddUserToProject(t, testApp, project.ID, student.ID, "team_developer")

	// Create some data to test with
	userStory := CreateTestUserStory(t, testApp, "Dashboard US", project.ID)

	// Student's task
	task := CreateTestTask(t, testApp, "Student Task 1", userStory.ID, student.ID)

	// Task for evaluation
	evalTask := CreateTestTask(t, testApp, "Task for Eval", userStory.ID, student.ID)
	evalTask.SubmittedForEvaluation = true
	now := time.Now()
	evalTask.SubmittedAt = &now
	evalTask.SubmittedByID = &student.ID
	err = testApp.DB.Save(evalTask).Error
	require.NoError(t, err)


	t.Run("GET /api/users/:id/tasks - Student can get their own tasks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d/tasks", student.ID), nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
		rec := httptest.NewRecorder()
		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var body map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &body)
		require.NoError(t, err)

		tasks, ok := body["tasks"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, tasks)

		firstTask := tasks[0].(map[string]interface{})
		assert.Contains(t, firstTask["title"], "Task")
	})

	t.Run("GET /api/users/:id/tasks - Student CANNOT get another student's tasks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%d/tasks", student2.ID), nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
		rec := httptest.NewRecorder()
		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("GET /api/evaluations/pending - Admin (as Teacher) can get pending evaluations", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/evaluations/pending", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rec := httptest.NewRecorder()
		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var body map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &body)
		require.NoError(t, err)

		evaluations, ok := body["evaluations"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, evaluations)
	})

	t.Run("GET /api/evaluations/pending - Student CANNOT get pending evaluations", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/evaluations/pending", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
		rec := httptest.NewRecorder()
		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("GET /api/activities/recent - Student can get recent activities", func(t *testing.T) {
		// Create an activity by completing a task
		_, err := testApp.TaskService.UpdateTaskStatus(task.ID, string(models.StatusDone), student.ID)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/activities/recent", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
		rec := httptest.NewRecorder()
		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var body map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &body)
		require.NoError(t, err)

		activities, ok := body["activities"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, activities, "Expected to find at least one activity")

		firstActivity := activities[0].(map[string]interface{})
		assert.Equal(t, "task_completed", firstActivity["type"])
	})

	t.Run("GET /api/deadlines/upcoming - Student can get upcoming deadlines", func(t *testing.T) {
		// Create a deadline for the project
		deadlineDate := time.Now().Add(10 * 24 * time.Hour)
		deadline := &models.Deadline{
			Title:     "Project Final Deadline",
			Type:      "project_milestone",
			ProjectID: &project.ID,
			Date:      deadlineDate,
		}
		err := testApp.DB.Create(deadline).Error
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/deadlines/upcoming?days=15", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
		rec := httptest.NewRecorder()
		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var body map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &body)
		require.NoError(t, err)

		deadlines, ok := body["deadlines"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, deadlines)

		firstDeadline := deadlines[0].(map[string]interface{})
		assert.Equal(t, "Project Final Deadline", firstDeadline["title"])
		assert.Equal(t, "project_milestone", firstDeadline["type"])
	})

}
