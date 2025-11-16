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

	t.Run("Get Sprints with Stats", func(t *testing.T) {
		// Create a new project for this test
		p := CreateTestProject(t, testApp, "Sprint Stats Project", creator.ID)

		// Create a sprint
		s := &models.Sprint{Name: "Stats Sprint"}
		err := sprintService.CreateSprint(s, p.ID, creator.ID)
		require.NoError(t, err)

		// Create user stories and tasks
		us := CreateTestUserStory(t, testApp, "Stats US", p.ID)
		us.SprintID = &s.ID // Assign to sprint
		err = testApp.DB.Save(us).Error
		require.NoError(t, err)

		task1 := CreateTestTask(t, testApp, "Done Task", us.ID, creator.ID)
		task1.Status = models.StatusDone
		err = testApp.DB.Save(task1).Error
		require.NoError(t, err)

		task2 := CreateTestTask(t, testApp, "In Progress Task", us.ID, creator.ID)
		task2.Status = models.StatusInProgress
		err = testApp.DB.Save(task2).Error
		require.NoError(t, err)

		_ = CreateTestTask(t, testApp, "Todo Task", us.ID, creator.ID) // Status is 'todo' by default

		// Call the service method
		sprintsWithStats, _, err := sprintService.GetSprintsByProjectID(p.ID)
		require.NoError(t, err)

		// Assertions
		require.Len(t, sprintsWithStats, 1)
		stats := sprintsWithStats[0]

		assert.Equal(t, 3, stats.TotalTasks)
		assert.Equal(t, 1, stats.CompletedTasks)
		assert.Equal(t, 1, stats.InProgressTasks)
		assert.Equal(t, 1, stats.PendingTasks)
		assert.InDelta(t, 33.33, stats.Progress, 0.01)
	})
}
