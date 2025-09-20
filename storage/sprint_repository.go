package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// SprintRepository handles database operations for sprints.
type SprintRepository struct {
	DB *gorm.DB
}

// NewSprintRepository creates a new instance of SprintRepository.
func NewSprintRepository(db *gorm.DB) *SprintRepository {
	return &SprintRepository{DB: db}
}

// CreateSprint adds a new sprint to the database.
func (r *SprintRepository) CreateSprint(sprint *models.Sprint) error {
	return r.DB.Create(sprint).Error
}

// GetSprintsByProjectID retrieves all sprints for a given project ID.
func (r *SprintRepository) GetSprintsByProjectID(projectID uint) ([]models.Sprint, error) {
	var sprints []models.Sprint
	err := r.DB.Where("project_id = ?", projectID).Preload("CreatedBy").Find(&sprints).Error
	return sprints, err
}

// GetSprintByID retrieves a single sprint by its ID.
func (r *SprintRepository) GetSprintByID(id uint) (*models.Sprint, error) {
	var sprint models.Sprint
	err := r.DB.Preload("CreatedBy").Preload("Project").First(&sprint, id).Error
	return &sprint, err
}

// UpdateSprint updates an existing sprint in the database.
func (r *SprintRepository) UpdateSprint(sprint *models.Sprint) error {
	return r.DB.Save(sprint).Error
}

// DeleteSprint removes a sprint from the database by its ID.
func (r *SprintRepository) DeleteSprint(id uint) error {
	return r.DB.Delete(&models.Sprint{}, id).Error
}

// DeleteSprintsByProjectID deletes all sprints associated with a project.
func (r *SprintRepository) DeleteSprintsByProjectID(tx *gorm.DB, projectID uint) error {
	return tx.Where("project_id = ?", projectID).Delete(&models.Sprint{}).Error
}
