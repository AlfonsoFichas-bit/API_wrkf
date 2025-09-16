package services

import (
	"API_wrkf/models"
	"API_wrkf/storage"
	"fmt"

	"gorm.io/gorm"
)

// ProjectService handles the business logic for projects.
type ProjectService struct {
	Repo *storage.ProjectRepository
}

// NewProjectService creates a new instance of ProjectService.
func NewProjectService(repo *storage.ProjectRepository) *ProjectService {
	return &ProjectService{Repo: repo}
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

	return member, nil
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

	existingProject.Name = updates["Name"].(string)
	existingProject.Description = updates["Description"].(string)

	if err := s.Repo.UpdateProject(existingProject); err != nil {
		return nil, err
	}

	return existingProject, nil
}

// DeleteProject handles the transactional deletion of a project and its members.
func (s *ProjectService) DeleteProject(projectID uint, requestingUserID uint, requestingUserRole string) error {
	// 1. Get the existing project to check for ownership.
	existingProject, err := s.Repo.GetProjectByID(projectID)
	if err != nil {
		return fmt.Errorf("project not found")
	}

	// 2. Check permissions.
	isOwner := existingProject.CreatedByID == requestingUserID
	isAdmin := requestingUserRole == string(models.RoleAdmin)
	if !isOwner && !isAdmin {
		return fmt.Errorf("forbidden: you do not have permission to delete this project")
	}

	// 3. Perform deletion within a transaction.
	return s.Repo.DB.Transaction(func(tx *gorm.DB) error {
		// First, delete all members associated with the project.
		if err := s.Repo.DeleteProjectMembersByProjectID(tx, projectID); err != nil {
			return err // Rollback
		}

		// Then, delete the project itself.
		if err := s.Repo.DeleteProject(tx, projectID); err != nil {
			return err // Rollback
		}

		// Return nil to commit the transaction
		return nil
	})
}
