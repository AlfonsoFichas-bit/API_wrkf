package tests

import (
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectService(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	projectService := testApp.ProjectService

	// Create a test user to be the project creator
	creator, _ := CreateTestUser(t, testApp, "project_creator@test.com", "admin")

	t.Run("Create Project", func(t *testing.T) {
		projectName := "Test Project Alpha"
		project := &models.Project{
			Name:        projectName,
			Description: "A test project",
		}
		err := projectService.CreateProject(project, creator.ID)
		require.NoError(t, err)
		assert.Equal(t, projectName, project.Name)
		assert.Equal(t, creator.ID, project.CreatedByID)

		// Verify that the creator was added as a member
		members, err := projectService.GetProjectMembers(project.ID)
		require.NoError(t, err)
		assert.Len(t, members, 1)
		assert.Equal(t, creator.ID, members[0].UserID)
		assert.Equal(t, "product_owner", members[0].Role)
	})

	t.Run("Add and Get Members", func(t *testing.T) {
		project := &models.Project{Name: "Membership Test Project"}
		err := projectService.CreateProject(project, creator.ID)
		require.NoError(t, err)

		// Create another user to add to the project
		newUser, _ := CreateTestUser(t, testApp, "new_member@test.com", "user")
		_, err = projectService.AddMemberToProject(project.ID, newUser.ID, "team_developer")
		require.NoError(t, err)

		members, err := projectService.GetProjectMembers(project.ID)
		require.NoError(t, err)
		assert.Len(t, members, 2) // Creator + new member
	})

	t.Run("Get Unassigned Users", func(t *testing.T) {
		project := &models.Project{Name: "Unassigned Users Test"}
		err := projectService.CreateProject(project, creator.ID)
		require.NoError(t, err)

		// Create some users, one of whom will be in the project
		userA, _ := CreateTestUser(t, testApp, "userA@test.com", "user")
		userB, _ := CreateTestUser(t, testApp, "userB@test.com", "user")
		_, err = projectService.AddMemberToProject(project.ID, userA.ID, "team_developer")
		require.NoError(t, err)

		unassigned, err := projectService.GetUnassignedUsers(project.ID)
		require.NoError(t, err)

		// Check that userB is in the unassigned list
		foundUserB := false
		for _, u := range unassigned {
			if u.ID == userB.ID {
				foundUserB = true
				break
			}
		}
		assert.True(t, foundUserB, "UserB should be in the unassigned list")

		// Check that userA is NOT in the unassigned list
		foundUserA := false
		for _, u := range unassigned {
			if u.ID == userA.ID {
				foundUserA = true
				break
			}
		}
		assert.False(t, foundUserA, "UserA should not be in the unassigned list")
	})
}
