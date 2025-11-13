package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// EvaluationRepository handles database operations for evaluations.
type EvaluationRepository struct {
	db *gorm.DB
}

// NewEvaluationRepository creates a new EvaluationRepository.
func NewEvaluationRepository(db *gorm.DB) *EvaluationRepository {
	return &EvaluationRepository{db: db}
}

// CreateEvaluation creates a new evaluation in the database.
// It uses a transaction to ensure that the evaluation and all its criterion evaluations are created atomically.
func (r *EvaluationRepository) CreateEvaluation(evaluation *models.Evaluation) error {
	return r.db.Create(evaluation).Error
}

// GetEvaluationsByTaskID retrieves all evaluations for a given task,
// preloading related data for a complete view.
func (r *EvaluationRepository) GetEvaluationsByTaskID(taskID uint) ([]models.Evaluation, error) {
	var evaluations []models.Evaluation
	err := r.db.
		Preload("Evaluator").
		Preload("Rubric").
		Preload("CriterionEvaluations").
		Preload("CriterionEvaluations.Criterion").
		Where("task_id = ?", taskID).
		Find(&evaluations).Error
	return evaluations, err
}
