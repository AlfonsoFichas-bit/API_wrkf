package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// EvaluationRepository handles database operations for evaluations.
type EvaluationRepository struct {
	DB *gorm.DB
}

// NewEvaluationRepository creates a new instance of EvaluationRepository.
func NewEvaluationRepository(db *gorm.DB) *EvaluationRepository {
	return &EvaluationRepository{DB: db}
}

// CreateEvaluation inserts a new evaluation record into the database.
func (r *EvaluationRepository) CreateEvaluation(evaluation *models.Evaluation) error {
	return r.DB.Create(evaluation).Error
}

// GetEvaluationByTaskID retrieves an evaluation for a specific task, preloading related data.
func (r *EvaluationRepository) GetEvaluationByTaskID(taskID uint) (*models.Evaluation, error) {
	var evaluation models.Evaluation
	err := r.DB.Where("task_id = ?", taskID).
		Preload("Evaluator").
		Preload("Task").
		Preload("Criteria.RubricCriterion").
		Preload("Feedback.Author").
		First(&evaluation).Error
	return &evaluation, err
}

// UpdateEvaluation updates an existing evaluation record in the database.
func (r *EvaluationRepository) UpdateEvaluation(evaluation *models.Evaluation) error {
	return r.DB.Save(evaluation).Error
}
