
package services

import (
	"fmt"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"

	"gorm.io/gorm"
)

// ProjectService handles the business logic for projects.
type ProjectService struct {
	Repo            *storage.ProjectRepository
	UserStoryRepo   *storage.UserStoryRepository // NEW
	SprintRepo      *storage.SprintRepository    // NEW
	TaskRepo        *storage.TaskRepository      // NEW
}

// NewProjectService creates a new instance of ProjectService.
func NewProjectService(repo *storage.ProjectRepository, userStoryRepo *storage.UserStoryRepository, sprintRepo *storage.SprintRepository, taskRepo *storage.TaskRepository) *ProjectService {
	return &ProjectService{
		Repo:            repo,
		UserStoryRepo:   userStoryRepo,
		SprintRepo:      sprintRepo,
		TaskRepo:        taskRepo,
	}
}

// CreateProject handles the business logic for creating a new project.
func (s *ProjectService) CreateProject(project *models.Project, creatorID uint) error {
	project.CreatedByID = creatorID
	return s.Repo.CreateProject(project)
}

// AddMemberToProject handles the business logic for adding a user to a project.
func (s *ProjectService) AddMemberToProject(projectID, userID uint, role string) (*models.ProjectMember, error) {
	projectRole := models.ProjectRole(role)
	if !projectRole.IsValid() {
		return nil, fmt.Errorf("invalid project role: '%s'", role)
	}

	member := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}

	if err := s.Repo.AddMemberToProject(member); err != nil {
		return nil, err
	}

	return s.Repo.GetProjectMemberByID(member.ID)
}

// GetProjects retrieves all projects.
func (s *ProjectService) GetProjects() ([]models.Project, error) {
	return s.Repo.GetProjects()
}

// GetProjectByID retrieves a single project by its ID.
func (s *ProjectService) GetProjectByID(id uint) (*models.Project, error) {
	return s.Repo.GetProjectByID(id)
}

// UpdateProject handles the business logic for updating a project, including permission checks.
func (s *ProjectService) UpdateProject(projectID uint, updates map[string]interface{}, requestingUserID uint, requestingUserRole string) (*models.Project, error) {
	existingProject, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	isOwner := existingProject.CreatedByID == requestingUserID
	isAdmin := requestingUserRole == string(models.RoleAdmin)
	if !isOwner && !isAdmin {
		return nil, fmt.Errorf("forbidden: you do not have permission to update this project")
	}

	if name, ok := updates["Name"].(string); ok {
		existingProject.Name = name
	}
	if description, ok := updates["Description"].(string); ok {
		existingProject.Description = description
	}

	if err := s.Repo.UpdateProject(existingProject); err != nil {
		return nil, err
	}

	return existingProject, nil
}

// DeleteProject handles the transactional deletion of a project and all its dependencies.
func (s *ProjectService) DeleteProject(projectID uint, requestingUserID uint, requestingUserRole string) error {
	existingProject, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		return fmt.Errorf("project not found")
	}

	isOwner := existingProject.CreatedByID == requestingUserID
	isAdmin := requestingUserRole == string(models.RoleAdmin)
	if !isOwner && !isAdmin {
		return fmt.Errorf("forbidden: you do not have permission to delete this project")
	}

	// Perform all deletions within a single transaction.
	return s.Repo.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Get all User Story IDs for the project.
		storyIDs, err := s.UserStoryRepo.GetUserStoryIDsByProjectID(tx, projectID)
		if err != nil {
			return err // Rollback
		}

		// 2. If there are user stories, delete their dependent tasks first.
		if len(storyIDs) > 0 {
			if err := s.TaskRepo.DeleteTasksByUserStoryIDs(tx, storyIDs); err != nil {
				return err // Rollback
			}
		}

		// 3. Delete all User Stories for the project.
		if err := s.UserStoryRepo.DeleteUserStoriesByProjectID(tx, projectID); err != nil {
			return err // Rollback
		}

		// 4. Delete all Sprints for the project.
		if err := s.SprintRepo.DeleteSprintsByProjectID(tx, projectID); err != nil {
			return err // Rollback
		}

		// 5. Delete all Project Members for the project.
		if err := s.Repo.DeleteProjectMembersByProjectID(tx, projectID); err != nil {
			return err // Rollback
		}

		// 6. Finally, delete the project itself.
		if err := s.Repo.DeleteProject(tx, projectID); err != nil {
			return err // Rollback
		}

		return nil // Commit
	})
}

// GetUserRoleInProject retrieves a user's role within a specific project.
func (s *ProjectService) GetUserRoleInProject(userID, projectID uint) (string, error) {
	return s.Repo.GetUserRoleInProject(userID, projectID)
}
