package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// IRubricRepository defines the interface for rubric data operations.
type IRubricRepository interface {
	CreateRubric(rubric *models.Rubric) error
	GetAllRubrics(filters map[string]interface{}) ([]models.Rubric, error)
	GetRubricByID(id uint) (*models.Rubric, error)
	UpdateRubric(rubric *models.Rubric) error
	DeleteRubric(id uint) error
}

// RubricRepository implements the IRubricRepository interface.
type RubricRepository struct {
	db *gorm.DB
}

// NewRubricRepository creates a new instance of RubricRepository.
func NewRubricRepository(db *gorm.DB) IRubricRepository {
	return &RubricRepository{db: db}
}

// CreateRubric saves a new rubric to the database.
func (r *RubricRepository) CreateRubric(rubric *models.Rubric) error {
	return r.db.Create(rubric).Error
}

// GetAllRubrics retrieves all rubrics from the database, with optional filters.
func (r *RubricRepository) GetAllRubrics(filters map[string]interface{}) ([]models.Rubric, error) {
	var rubrics []models.Rubric
	query := r.db.Model(&models.Rubric{})

	if len(filters) > 0 {
		for key, value := range filters {
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Preload("Criteria.Levels").Find(&rubrics).Error
	return rubrics, err
}

// GetRubricByID retrieves a single rubric by its ID.
func (r *RubricRepository) GetRubricByID(id uint) (*models.Rubric, error) {
	var rubric models.Rubric
	err := r.db.Preload("Criteria.Levels").First(&rubric, id).Error
	return &rubric, err
}

// UpdateRubric modifies an existing rubric in the database.
func (r *RubricRepository) UpdateRubric(rubric *models.Rubric) error {
	return r.db.Save(rubric).Error
}

// DeleteRubric removes a rubric from the database.
func (r *RubricRepository) DeleteRubric(id uint) error {
	return r.db.Delete(&models.Rubric{}, id).Error
}
