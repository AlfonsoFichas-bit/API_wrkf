
package services

import (
	"fmt"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// UserStoryService handles the business logic for user stories.
type UserStoryService struct {
	Repo           *storage.UserStoryRepository
	ProjectService *ProjectService // Dependency to check project-level roles
}

// NewUserStoryService creates a new instance of UserStoryService.
func NewUserStoryService(repo *storage.UserStoryRepository, projectService *ProjectService) *UserStoryService {
	return &UserStoryService{
		Repo:           repo,
		ProjectService: projectService,
	}
}

// CreateUserStory handles the business logic for creating a new user story.
func (s *UserStoryService) CreateUserStory(userStory *models.UserStory, projectID uint, creatorID uint) error {
	userStory.ProjectID = projectID
	userStory.CreatedByID = creatorID
	return s.Repo.CreateUserStory(userStory)
}

// GetUserStoriesByProjectID retrieves all user stories for a specific project.
func (s *UserStoryService) GetUserStoriesByProjectID(projectID uint) ([]models.UserStory, error) {
	return s.Repo.GetUserStoriesByProjectID(projectID)
}

// checkUserStoryPermissions is a helper function to verify if a user can modify a user story.
func (s *UserStoryService) checkUserStoryPermissions(userStoryID, requestingUserID uint, platformRole string) (bool, error) {
	// First, check if the user is a platform admin.
	if platformRole == string(models.RoleAdmin) {
		return true, nil
	}

	// Get the user story to find out which project it belongs to.
	userStory, err := s.Repo.GetUserStoryByID(userStoryID)
	if err != nil {
		return false, fmt.Errorf("user story not found")
	}

	// Get the user's role within that specific project.
	projectRole, err := s.ProjectService.GetUserRoleInProject(requestingUserID, userStory.ProjectID)
	if err != nil {
		// This includes the case where the user is not a member of the project.
		return false, nil
	}

	// Check if the project role has modification permissions.
	switch models.ProjectRole(projectRole) {
	case models.RoleProductOwner, models.RoleScrumMaster:
		return true, nil
	default:
		return false, nil
	}
}

// UpdateUserStory handles updating a user story after checking permissions.
func (s *UserStoryService) UpdateUserStory(storyID, requestingUserID uint, platformRole string, updates map[string]interface{}) (*models.UserStory, error) {
	canUpdate, err := s.checkUserStoryPermissions(storyID, requestingUserID, platformRole)
	if err != nil {
		return nil, err
	}
	if !canUpdate {
		return nil, fmt.Errorf("forbidden: you do not have permission to update this user story")
	}

	existingStory, err := s.Repo.GetUserStoryByID(storyID)
	if err != nil {
		return nil, fmt.Errorf("user story not found")
	}

	// Apply updates (a more robust implementation would use reflection or a library)
	if title, ok := updates["title"].(string); ok {
		existingStory.Title = title
	}
	if description, ok := updates["description"].(string); ok {
		existingStory.Description = description
	}

	if err := s.Repo.UpdateUserStory(existingStory); err != nil {
		return nil, err
	}
	return existingStory, nil
}

// DeleteUserStory handles deleting a user story after checking permissions.
func (s *UserStoryService) DeleteUserStory(storyID, requestingUserID uint, platformRole string) error {
	canDelete, err := s.checkUserStoryPermissions(storyID, requestingUserID, platformRole)
	if err != nil {
		return err
	}
	if !canDelete {
		return fmt.Errorf("forbidden: you do not have permission to delete this user story")
	}

	return s.Repo.DeleteUserStory(storyID)
}
