package tests

import (
	"testing"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSprintService(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	sprintService := testApp.SprintService

	creator, _ := CreateTestUser(t, testApp, "sprint_creator@test.com", "admin")
	project := CreateTestProject(t, testApp, "Sprint Test Project", creator.ID)

	t.Run("Create and Get Sprint", func(t *testing.T) {
		startDate := time.Now().Add(24 * time.Hour)
		endDate := startDate.Add(14 * 24 * time.Hour)
		sprint := &models.Sprint{
			Name:      "Test Sprint 1",
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		err := sprintService.CreateSprint(sprint, project.ID, creator.ID)
		require.NoError(t, err)
		assert.NotZero(t, sprint.ID)

		found, err := sprintService.GetSprintByID(sprint.ID)
		require.NoError(t, err)
		assert.Equal(t, "Test Sprint 1", found.Name)
	})

	t.Run("Update Sprint Status", func(t *testing.T) {
		sprint := &models.Sprint{Name: "Status Update Sprint"}
		err := sprintService.CreateSprint(sprint, project.ID, creator.ID)
		require.NoError(t, err)

		// Start the sprint
		err = sprintService.UpdateSprintStatus(sprint.ID, "active")
		require.NoError(t, err)

		updated, err := sprintService.GetSprintByID(sprint.ID)
		require.NoError(t, err)
		assert.Equal(t, "active", updated.Status)

		// Try to start another sprint (should fail)
		sprint2 := &models.Sprint{Name: "Another Sprint"}
		err = sprintService.CreateSprint(sprint2, project.ID, creator.ID)
		require.NoError(t, err)
		err = sprintService.UpdateSprintStatus(sprint2.ID, "active")
		assert.Error(t, err, "Should not be able to start a new sprint while one is active")

		// End the first sprint
		err = sprintService.UpdateSprintStatus(sprint.ID, "completed")
		require.NoError(t, err)

		updated, err = sprintService.GetSprintByID(sprint.ID)
		require.NoError(t, err)
		assert.Equal(t, "completed", updated.Status)
	})
}
