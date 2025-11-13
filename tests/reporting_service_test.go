package tests

import (
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

func TestReportingEndpoints(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	creator, creatorToken := CreateTestUser(t, testApp, "report_user@test.com", "admin")
	project := CreateTestProject(t, testApp, "Reporting Project", creator.ID)

	t.Run("Sprint Commitment Report", func(t *testing.T) {
		// --- Setup Data ---
		sprint := &models.Sprint{
			Name:      "Commitment Sprint",
			ProjectID: project.ID,
		}
		err := testApp.DB.Create(sprint).Error
		require.NoError(t, err)

		points5 := 5
		points8 := 8
		points3 := 3

		// Story 1: Done (5 points)
		us1 := &models.UserStory{Title: "US1", ProjectID: project.ID, SprintID: &sprint.ID, Points: &points5, Status: "done"}
		// Story 2: In Progress (8 points)
		us2 := &models.UserStory{Title: "US2", ProjectID: project.ID, SprintID: &sprint.ID, Points: &points8, Status: "in_progress"}
		// Story 3: Done (3 points)
		us3 := &models.UserStory{Title: "US3", ProjectID: project.ID, SprintID: &sprint.ID, Points: &points3, Status: "done"}

		require.NoError(t, testApp.DB.Create(&us1).Error)
		require.NoError(t, testApp.DB.Create(&us2).Error)
		require.NoError(t, testApp.DB.Create(&us3).Error)

		// --- Make API Call ---
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/sprints/%d/reports/commitment", sprint.ID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+creatorToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		// --- Assertions ---
		assert.Equal(t, http.StatusOK, rec.Code)

		var report models.CommitmentReport
		err = json.Unmarshal(rec.Body.Bytes(), &report)
		require.NoError(t, err)

		assert.Equal(t, sprint.ID, report.SprintID)
		assert.Equal(t, "Commitment Sprint", report.SprintName)
		assert.Equal(t, 16, report.CommittedPoints) // 5 + 8 + 3
		assert.Equal(t, 8, report.CompletedPoints)  // 5 + 3
		assert.InDelta(t, 50.0, report.CompletionRate, 0.01) // (8 / 16) * 100
	})

	t.Run("Sprint Burndown Chart", func(t *testing.T) {
		// --- Setup Data ---
		startDate := time.Now()
		endDate := startDate.Add(4 * 24 * time.Hour) // 5-day sprint
		sprint := &models.Sprint{
			Name:      "Burndown Sprint",
			ProjectID: project.ID,
			StartDate: &startDate,
			EndDate:   &endDate,
		}
		err := testApp.DB.Create(sprint).Error
		require.NoError(t, err)

		points8 := 8
		us := &models.UserStory{
			Title:     "Burndown US",
			ProjectID: project.ID,
			SprintID:  &sprint.ID,
			Points:    &points8,
			Status:    "in_progress",
		}
		require.NoError(t, testApp.DB.Create(us).Error)

		task1 := &models.Task{Title: "Task 1", UserStoryID: us.ID, Status: "done"}
		task2 := &models.Task{Title: "Task 2", UserStoryID: us.ID, Status: "done"}
		require.NoError(t, testApp.DB.Create(task1).Error)
		require.NoError(t, testApp.DB.Create(task2).Error)

		// Simulate task completion history
		day2Completion := startDate.Add(1 * 24 * time.Hour) // Completed on Day 2
		history1 := &models.TaskHistory{TaskID: task1.ID, FieldName: "status", OldValue: "in_progress", NewValue: "done", ChangedAt: day2Completion}
		history2 := &models.TaskHistory{TaskID: task2.ID, FieldName: "status", OldValue: "in_progress", NewValue: "done", ChangedAt: day2Completion}
		require.NoError(t, testApp.DB.Create(history1).Error)
		require.NoError(t, testApp.DB.Create(history2).Error)

		// --- Make API Call ---
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/sprints/%d/reports/burndown", sprint.ID), nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+creatorToken)
		rec := httptest.NewRecorder()

		testApp.Router.ServeHTTP(rec, req)

		// --- Assertions ---
		assert.Equal(t, http.StatusOK, rec.Code, "API call should be successful")

		var report models.BurndownReport
		err = json.Unmarshal(rec.Body.Bytes(), &report)
		require.NoError(t, err, "Should be able to unmarshal the response")

		assert.Equal(t, sprint.ID, report.SprintID)
		assert.Equal(t, 8, report.TotalPoints)
		require.Len(t, report.BurndownData, 5, "Burndown should have data for 5 days")

		// Ideal line should burn 2 points per day (8 points / 4 intervals)
		assert.InDelta(t, 8.0, report.BurndownData[0].IdealPoints, 0.01)
		assert.InDelta(t, 6.0, report.BurndownData[1].IdealPoints, 0.01)
		assert.InDelta(t, 4.0, report.BurndownData[2].IdealPoints, 0.01)
		assert.InDelta(t, 2.0, report.BurndownData[3].IdealPoints, 0.01)
		assert.InDelta(t, 0.0, report.BurndownData[4].IdealPoints, 0.01)

		// Actual line should remain at 8, then drop to 0 on day 2
		assert.InDelta(t, 8.0, report.BurndownData[0].RemainingPoints, 0.01, "Day 1 should have 8 points remaining")
		assert.InDelta(t, 0.0, report.BurndownData[1].RemainingPoints, 0.01, "Day 2 should have 0 points remaining")
		assert.InDelta(t, 0.0, report.BurndownData[2].RemainingPoints, 0.01, "Day 3 should have 0 points remaining")
		assert.InDelta(t, 0.0, report.BurndownData[3].RemainingPoints, 0.01, "Day 4 should have 0 points remaining")
		assert.InDelta(t, 0.0, report.BurndownData[4].RemainingPoints, 0.01, "Day 5 should have 0 points remaining")
	})
}
