package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventEndpoints(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	// --- Setup Data ---
	member, memberToken := CreateTestUser(t, testApp, "event_member@test.com", "user")
	_, nonMemberToken := CreateTestUser(t, testApp, "event_nonmember@test.com", "user")
	project := CreateTestProject(t, testApp, "Event Project", member.ID)
	AddUserToProject(t, testApp, project.ID, member.ID, "team_developer")

	var createdEvent models.Event

	// --- Test Cases ---
	t.Run("Non-member cannot create event", func(t *testing.T) {
		event := map[string]string{
			"title":     "Unauthorized Event",
			"startDate": time.Now().Format(time.RFC3339),
			"endDate":   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}
		reqBody, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/projects/%d/events", project.ID), bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+nonMemberToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Should fail with a permission error")
	})

	t.Run("Member can create event", func(t *testing.T) {
		startDate := time.Now().Add(24 * time.Hour)
		endDate := startDate.Add(1 * time.Hour)
		event := map[string]string{
			"title":     "Project Kick-off Meeting",
			"startDate": startDate.Format(time.RFC3339),
			"endDate":   endDate.Format(time.RFC3339),
		}
		reqBody, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/projects/%d/events", project.ID), bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+memberToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		err := json.Unmarshal(rec.Body.Bytes(), &createdEvent)
		require.NoError(t, err)
		assert.Equal(t, "Project Kick-off Meeting", createdEvent.Title)
		assert.Equal(t, project.ID, createdEvent.ProjectID)
		assert.Equal(t, member.ID, createdEvent.CreatedByID)
	})

	t.Run("Get events for project", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/projects/%d/events", project.ID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+memberToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var events []models.Event
		err := json.Unmarshal(rec.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 1)
		assert.Equal(t, createdEvent.ID, events[0].ID)
	})

	t.Run("Update event", func(t *testing.T) {
		updates := map[string]string{"title": "Updated Meeting Title"}
		reqBody, _ := json.Marshal(updates)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/events/%d", createdEvent.ID), bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+memberToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var updatedEvent models.Event
		err := json.Unmarshal(rec.Body.Bytes(), &updatedEvent)
		require.NoError(t, err)
		assert.Equal(t, "Updated Meeting Title", updatedEvent.Title)
	})

	t.Run("Delete event", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/events/%d", createdEvent.ID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+memberToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// Verify it's gone
		reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/projects/%d/events", project.ID), nil)
		reqGet.Header.Set(echo.HeaderAuthorization, "Bearer "+memberToken)
		recGet := httptest.NewRecorder()
		testApp.Router.ServeHTTP(recGet, reqGet)
		var events []models.Event
		json.Unmarshal(recGet.Body.Bytes(), &events)
		assert.Len(t, events, 0)
	})
}
