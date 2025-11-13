package tests

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportEndpoint(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// --- Setup Data ---
	user, userToken := CreateTestUser(t, testApp, "export_user@test.com", "user")
	project := CreateTestProject(t, testApp, "Export Test Project", user.ID)
	us1 := CreateTestUserStory(t, testApp, "Export US 1", project.ID)
	task1 := CreateTestTask(t, testApp, "Export Task 1", us1.ID, user.ID)

	// Create a user story with no tasks to test that case as well
	us2 := CreateTestUserStory(t, testApp, "Export US 2 (No Tasks)", project.ID)

	// --- Make API Call ---
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/projects/%d/export", project.ID), nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+userToken)
	rec := httptest.NewRecorder()

	testApp.Router.ServeHTTP(rec, req)

	// --- Assertions ---
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check headers
	assert.Equal(t, "text/csv", rec.Header().Get("Content-Type"))
	assert.True(t, strings.HasPrefix(rec.Header().Get("Content-Disposition"), "attachment; filename="), "Content-Disposition should be an attachment")

	// Check CSV content
	reader := csv.NewReader(rec.Body)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Expected: 1 header row + 1 row for task1 + 1 row for us2 (with no task) = 3 total rows
	require.Len(t, records, 3, "CSV should have a header row and two data rows")

	// Check header
	expectedHeader := []string{
		"Project Name", "User Story ID", "User Story Title", "User Story Status",
		"Task ID", "Task Title", "Task Status", "Assigned To",
	}
	assert.Equal(t, expectedHeader, records[0])

	// Check data row for task1
	assert.Equal(t, project.Name, records[1][0])
	assert.Equal(t, fmt.Sprintf("%d", us1.ID), records[1][1])
	assert.Equal(t, us1.Title, records[1][2])
	assert.Equal(t, fmt.Sprintf("%d", task1.ID), records[1][4])
	assert.Equal(t, task1.Title, records[1][5])

	// Check data row for us2 (no tasks)
	assert.Equal(t, project.Name, records[2][0])
	assert.Equal(t, fmt.Sprintf("%d", us2.ID), records[2][1])
	assert.Equal(t, us2.Title, records[2][2])
	assert.Equal(t, "", records[2][4], "Task ID should be empty for user story with no tasks")
}
