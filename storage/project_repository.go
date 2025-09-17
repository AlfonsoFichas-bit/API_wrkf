package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// ProjectRepository handles database operations for projects.
type ProjectRepository struct {
	DB *gorm.DB
}

// NewProjectRepository creates a new instance of ProjectRepository.
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

// CreateProject adds a new project to the database.
func (r *ProjectRepository) CreateProject(project *models.Project) error {
	return r.DB.Create(project).Error
}

// AddMemberToProject creates a new project membership record.
func (r *ProjectRepository) AddMemberToProject(member *models.ProjectMember) error {
	return r.DB.Create(member).Error
}

// GetProjects retrieves a list of all projects from the database.
func (r *ProjectRepository) GetProjects() ([]models.Project, error) {
	var projects []models.Project
	err := r.DB.Preload("CreatedBy").Find(&projects).Error
	return projects, err
}

// GetProjectByID retrieves a single project by its ID, including its members and their user details.
func (r *ProjectRepository) GetProjectByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.DB.Preload("CreatedBy").Preload("Members.User").First(&project, id).Error
	return &project, err
}

// UpdateProject updates an existing project in the database.
func (r *ProjectRepository) UpdateProject(project *models.Project) error {
	return r.DB.Save(project).Error
}

// DeleteProjectMembersByProjectID removes all membership records for a given project.
func (r *ProjectRepository) DeleteProjectMembersByProjectID(tx *gorm.DB, projectID uint) error {
	return tx.Where("project_id = ?", projectID).Delete(&models.ProjectMember{}).Error
}

// DeleteProject removes a project from the database by its ID.
func (r *ProjectRepository) DeleteProject(tx *gorm.DB, projectID uint) error {
	return tx.Delete(&models.Project{}, projectID).Error
}

// GetUserRoleInProject finds a user's specific role within a single project.
func (r *ProjectRepository) GetUserRoleInProject(userID, projectID uint) (string, error) {
	var member models.ProjectMember
	err := r.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&member).Error
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

// GetProjectMemberByID retrieves a single project membership record, preloading the user.
func (r *ProjectRepository) GetProjectMemberByID(id uint) (*models.ProjectMember, error) {
	var member models.ProjectMember
	err := r.DB.Preload("User").First(&member, id).Error
	return &member, err
}
