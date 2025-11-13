package tests

import (
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStoryService(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	userStoryService := testApp.UserStoryService
	sprintService := testApp.SprintService

	creator, _ := CreateTestUser(t, testApp, "us_creator@test.com", "admin")
	project := CreateTestProject(t, testApp, "User Story Test Project", creator.ID)

	t.Run("Create and Get User Story", func(t *testing.T) {
		us := &models.UserStory{
			Title:       "Test US",
			Description: "A description",
			ProjectID:   project.ID,
		}
		err := userStoryService.CreateUserStory(us, project.ID, creator.ID)
		require.NoError(t, err)
		assert.NotZero(t, us.ID)

		found, err := userStoryService.GetUserStoryByID(us.ID)
		require.NoError(t, err)
		assert.Equal(t, "Test US", found.Title)
	})

	t.Run("Assign User Story to Sprint", func(t *testing.T) {
		us := &models.UserStory{
			Title:     "Assignable US",
			ProjectID: project.ID,
		}
		err := userStoryService.CreateUserStory(us, project.ID, creator.ID)
		require.NoError(t, err)

		sprint := &models.Sprint{
			Name:      "Target Sprint",
			ProjectID: project.ID,
		}
		err = sprintService.CreateSprint(sprint, project.ID, creator.ID)
		require.NoError(t, err)

		// The AssignUserStoryToSprint method requires user info for permission checks.
		_, err = userStoryService.AssignUserStoryToSprint(sprint.ID, us.ID, creator.ID, "admin")
		require.NoError(t, err)

		found, err := userStoryService.GetUserStoryByID(us.ID)
		require.NoError(t, err)
		require.NotNil(t, found.SprintID)
		assert.Equal(t, sprint.ID, *found.SprintID)
	})
}
