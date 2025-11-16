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
		members, _, err := projectService.GetProjectMembers(project.ID)
		require.NoError(t, err)
		assert.Len(t, members, 1)
		assert.Equal(t, creator.ID, members[0].UserID)
		assert.Equal(t, "product_owner", members[0].Role)
	})

	t.Run("Add Member without Team (Regression Check)", func(t *testing.T) {
		project := &models.Project{Name: "Membership Regression Test Project"}
		err := projectService.CreateProject(project, creator.ID)
		require.NoError(t, err)

		// Create another user to add to the project
		newUser, _ := CreateTestUser(t, testApp, "new_member_regr@test.com", "user")

		// Add the member WITHOUT specifying a team
		addedMember, err := projectService.AddMemberToProject(project.ID, newUser.ID, "team_developer")
		require.NoError(t, err)
		assert.NotNil(t, addedMember)
		assert.Equal(t, newUser.ID, addedMember.UserID)
		assert.Nil(t, addedMember.TeamID, "TeamID should be nil when a member is added without a team")

		// Verify by fetching the members again
		members, _, err := projectService.GetProjectMembers(project.ID) // Note: ignoring teams result
		require.NoError(t, err)
		assert.Len(t, members, 2) // Creator + new member

		// Find the newly added member and double-check their TeamID
		var foundMember *models.ProjectMember
		for i := range members {
			if members[i].UserID == newUser.ID {
				foundMember = &members[i]
				break
			}
		}
		require.NotNil(t, foundMember, "Newly added member should be found in project members list")
		assert.Nil(t, foundMember.TeamID, "Fetched member's TeamID should also be nil")
	})

	t.Run("Get Project Members with Team Info", func(t *testing.T) {
		project := &models.Project{Name: "Team Info Test Project"}
		err := projectService.CreateProject(project, creator.ID)
		require.NoError(t, err)

		// Create a team
		team := &models.Team{Name: "Test Team", ProjectID: project.ID}
		err = testApp.DB.Create(team).Error
		require.NoError(t, err)

		// Create a user and add them to the project AND the team
		teamMemberUser, _ := CreateTestUser(t, testApp, "teammember@test.com", "user")
		memberWithTeam := &models.ProjectMember{
			ProjectID: project.ID,
			UserID:    teamMemberUser.ID,
			Role:      "team_developer",
			TeamID:    &team.ID,
		}
		err = testApp.DB.Create(memberWithTeam).Error
		require.NoError(t, err)

		// Get members and teams
		members, teams, err := projectService.GetProjectMembers(project.ID)
		require.NoError(t, err)

		// Assertions
		assert.Len(t, members, 2, "Should be 2 members (creator + new member)")
		assert.Len(t, teams, 1, "Should be 1 team in the project")
		assert.Equal(t, "Test Team", teams[0].Name)

		// Find the member with the team and check details
		var foundMember *models.ProjectMember
		for i := range members {
			if members[i].UserID == teamMemberUser.ID {
				foundMember = &members[i]
				break
			}
		}
		require.NotNil(t, foundMember)
		require.NotNil(t, foundMember.TeamID)
		assert.Equal(t, team.ID, *foundMember.TeamID)
		require.NotNil(t, foundMember.Team)
		assert.Equal(t, "Test Team", foundMember.Team.Name)
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
