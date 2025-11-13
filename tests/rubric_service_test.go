package tests

import (
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRubricService(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	rubricService := testApp.RubricService

	// Create a user and a project for context
	creator, _ := CreateTestUser(t, testApp, "rubric_creator@test.com", "admin")
	project := CreateTestProject(t, testApp, "Rubric Test Project", creator.ID)

	t.Run("Create and Get Rubric", func(t *testing.T) {
		rubric := &models.Rubric{
			Name:        "Test Rubric",
			Description: "A rubric for testing.",
			ProjectID:   project.ID,
			CreatedByID: creator.ID,
			Status:      models.RubricStatusDraft,
			Criteria: []models.RubricCriterion{
				{Title: "C1", MaxPoints: 10},
				{Title: "C2", MaxPoints: 5},
			},
		}

		err := rubricService.CreateRubric(rubric)
		require.NoError(t, err)
		assert.NotZero(t, rubric.ID)

		// Get the rubric back
		found, err := rubricService.GetRubricByID(rubric.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, "Test Rubric", found.Name)
		assert.Len(t, found.Criteria, 2)
	})

	t.Run("Update Rubric", func(t *testing.T) {
		rubric := &models.Rubric{
			Name:        "Update Me Rubric",
			ProjectID:   project.ID,
			CreatedByID: creator.ID,
		}
		err := rubricService.CreateRubric(rubric)
		require.NoError(t, err)

		// Update the name
		rubric.Name = "Updated Rubric Name"
		err = rubricService.UpdateRubric(rubric)
		require.NoError(t, err)

		found, err := rubricService.GetRubricByID(rubric.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Rubric Name", found.Name)
	})

	t.Run("Delete Rubric", func(t *testing.T) {
		rubric := &models.Rubric{
			Name:        "Delete Me Rubric",
			ProjectID:   project.ID,
			CreatedByID: creator.ID,
		}
		err := rubricService.CreateRubric(rubric)
		require.NoError(t, err)

		// Delete it
		err = rubricService.DeleteRubric(rubric.ID)
		require.NoError(t, err)

		// Try to find it again
		_, err = rubricService.GetRubricByID(rubric.ID)
		assert.Error(t, err, "Expected an error when finding a deleted rubric")
	})
}
