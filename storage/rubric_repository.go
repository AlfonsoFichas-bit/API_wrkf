package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// RubricRepository defines the interface for interacting with rubric data.
type RubricRepository interface {
	Create(rubric *models.Rubric) error
	FindAll(filters map[string]interface{}) ([]models.Rubric, error)
	FindByID(id uint) (*models.Rubric, error)
	Update(rubric *models.Rubric) error
	Delete(id uint) error
}

type rubricRepository struct {
	db *gorm.DB
}

// NewRubricRepository creates a new instance of RubricRepository.
func NewRubricRepository(db *gorm.DB) RubricRepository {
	return &rubricRepository{db: db}
}

// Create adds a new rubric to the database.
func (r *rubricRepository) Create(rubric *models.Rubric) error {
	return r.db.Create(rubric).Error
}

// FindAll retrieves all rubrics from the database, with preloaded criteria and levels.
// Filters can be applied, e.g., map[string]interface{}{"is_template": true}
func (r *rubricRepository) FindAll(filters map[string]interface{}) ([]models.Rubric, error) {
	var rubrics []models.Rubric
	query := r.db.Preload("Criteria.Levels")

	if len(filters) > 0 {
		query = query.Where(filters)
	}

	err := query.Find(&rubrics).Error
	return rubrics, err
}

// FindByID retrieves a single rubric by its ID, with preloaded criteria and levels.
func (r *rubricRepository) FindByID(id uint) (*models.Rubric, error) {
	var rubric models.Rubric
	err := r.db.Preload("Criteria.Levels").First(&rubric, id).Error
	if err != nil {
		return nil, err
	}
	return &rubric, nil
}

// Update modifies an existing rubric in the database.
// It performs a full update of the rubric and its associations.
func (r *rubricRepository) Update(rubric *models.Rubric) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(rubric).Error
}

// Delete removes a rubric from the database.
// The `constraint:OnDelete:CASCADE` in the model should handle deleting associated items.
func (r *rubricRepository) Delete(id uint) error {
	return r.db.Delete(&models.Rubric{}, id).Error
}
