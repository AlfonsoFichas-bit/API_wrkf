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
	SprintService  *SprintService  // Dependency to check sprint details
}

// NewUserStoryService creates a new instance of UserStoryService.
func NewUserStoryService(repo *storage.UserStoryRepository, projectService *ProjectService, sprintService *SprintService) *UserStoryService {
	return &UserStoryService{
		Repo:           repo,
		ProjectService: projectService,
		SprintService:  sprintService,
	}
}

// CreateUserStory handles the business logic for creating a new user story.
func (s *UserStoryService) CreateUserStory(userStory *models.UserStory, projectID uint, creatorID uint) error {
	if s == nil || s.Repo == nil {
		return fmt.Errorf("internal: user story repository not initialized")
	}
	userStory.ProjectID = projectID
	userStory.CreatedByID = creatorID
	return s.Repo.CreateUserStory(userStory)
}

// GetUserStoriesByProjectID retrieves all user stories for a specific project.
func (s *UserStoryService) GetUserStoriesByProjectID(projectID uint) ([]models.UserStory, error) {
	return s.Repo.GetUserStoriesByProjectID(projectID)
}

// GetUserStoryByID retrieves a single user story and manually hydrates the Sprint relationship.
func (s *UserStoryService) GetUserStoryByID(id uint) (*models.UserStory, error) {
	// 1. Get the base user story object.
	userStory, err := s.Repo.GetUserStoryByID(id)
	if err != nil {
		return nil, err
	}

	// 2. Manually hydrate the Sprint if a SprintID exists.
	if userStory.SprintID != nil {
		sprint, err := s.SprintService.GetSprintByID(*userStory.SprintID)
		// **THE FIX IS HERE: If hydration fails, return the error.**
		if err != nil {
			// This will tell us WHY the sprint isn't loading.
			return nil, fmt.Errorf("failed to hydrate sprint with ID %d: %w", *userStory.SprintID, err)
		}
		userStory.Sprint = sprint
	}

	return userStory, nil
}

// checkUserStoryPermissions is a helper function to verify if a user can modify a user story.
func (s *UserStoryService) checkUserStoryPermissions(userStoryID, requestingUserID uint, platformRole string) (bool, error) {
	if platformRole == string(models.RoleAdmin) {
		return true, nil
	}

	userStory, err := s.Repo.GetUserStoryByID(userStoryID)
	if err != nil {
		return false, fmt.Errorf("user story not found")
	}

	projectRole, err := s.ProjectService.GetUserRoleInProject(requestingUserID, userStory.ProjectID)
	if err != nil {
		return false, nil
	}

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

	if title, ok := updates["Title"].(string); ok {
		existingStory.Title = title
	}
	if description, ok := updates["Description"].(string); ok {
		existingStory.Description = description
	}
	if criteria, ok := updates["AcceptanceCriteria"].(string); ok {
		existingStory.AcceptanceCriteria = criteria
	}
	if priority, ok := updates["Priority"].(string); ok {
		existingStory.Priority = priority
	}
	if status, ok := updates["Status"].(string); ok {
		existingStory.Status = status
	}

	// CORRECCIÓN: Manejar Points (puede venir como int o float64)
	if pointsValue, ok := updates["Points"]; ok {
		switch v := pointsValue.(type) {
		case int:
			existingStory.Points = &v
		case float64:
			p := int(v)
			existingStory.Points = &p
		case *int:
			existingStory.Points = v
		}
	}

	// CORRECCIÓN: Manejar AssignedToID (puede venir como uint o float64)
	if assignedValue, ok := updates["AssignedToID"]; ok {
		switch v := assignedValue.(type) {
		case uint:
			existingStory.AssignedToID = &v
		case float64:
			id := uint(v)
			existingStory.AssignedToID = &id
		case *uint:
			existingStory.AssignedToID = v
		}
	}

	if err := s.Repo.UpdateUserStory(existingStory); err != nil {
		return nil, err
	}

	// Re-fetch the user story to ensure all fields, especially pointers, are correctly hydrated.
	return s.Repo.GetUserStoryByID(storyID)
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

// AssignUserStoryToSprint handles assigning a user story to a sprint with permission checks.
func (s *UserStoryService) AssignUserStoryToSprint(sprintID, storyID, requestingUserID uint, platformRole string) (*models.UserStory, error) {
	canAssign, err := s.checkUserStoryPermissions(storyID, requestingUserID, platformRole)
	if err != nil {
		return nil, err
	}
	if !canAssign {
		return nil, fmt.Errorf("forbidden: you do not have permission to assign this user story")
	}

	userStory, err := s.Repo.GetUserStoryByID(storyID)
	if err != nil {
		return nil, fmt.Errorf("user story not found")
	}
	sprint, err := s.SprintService.GetSprintByID(sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found")
	}

	if userStory.ProjectID != sprint.ProjectID {
		return nil, fmt.Errorf("cross-project assignment forbidden: user story and sprint belong to different projects")
	}

	userStory.SprintID = &sprintID
	if err := s.Repo.UpdateUserStory(userStory); err != nil {
		return nil, err
	}

	return userStory, nil
}
