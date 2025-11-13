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

func TestTaskComments(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// --- Create Test Data ---
	creator, creatorToken := CreateTestUser(t, testApp, "comment_creator@test.com", "user")
	project := CreateTestProject(t, testApp, "Comment Project", creator.ID)
	userStory := CreateTestUserStory(t, testApp, "Comment Story", project.ID)
	task := CreateTestTask(t, testApp, "Comment Task", userStory.ID, creator.ID)

	// --- Test Case 1: Add a comment to a task ---
	commentContent := "This is a test comment."
	body := map[string]string{"content": commentContent}
	reqBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/%d/comments", task.ID), bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+creatorToken)
	rec := httptest.NewRecorder()

	testApp.Router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var createdComment models.TaskComment
	err := json.Unmarshal(rec.Body.Bytes(), &createdComment)
	require.NoError(t, err)
	assert.Equal(t, commentContent, createdComment.Content)
	assert.Equal(t, creator.ID, createdComment.AuthorID)
	assert.Equal(t, task.ID, createdComment.TaskID)

	// --- Test Case 2: Get comments for the task ---
	reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d/comments", task.ID), nil)
	reqGet.Header.Set(echo.HeaderAuthorization, "Bearer "+creatorToken)
	recGet := httptest.NewRecorder()

	testApp.Router.ServeHTTP(recGet, reqGet)

	assert.Equal(t, http.StatusOK, recGet.Code)
	var fetchedComments []models.TaskComment
	err = json.Unmarshal(recGet.Body.Bytes(), &fetchedComments)
	require.NoError(t, err)

	require.Len(t, fetchedComments, 1)
	assert.Equal(t, createdComment.ID, fetchedComments[0].ID)
	assert.Equal(t, commentContent, fetchedComments[0].Content)
	require.NotNil(t, fetchedComments[0].Author)
	assert.Equal(t, creator.Nombre, fetchedComments[0].Author.Nombre)
}
