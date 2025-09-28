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

// CreateEvaluation adds a new evaluation to the database.
func (r *EvaluationRepository) CreateEvaluation(evaluation *models.Evaluation) error {
	return r.DB.Create(evaluation).Error
}

// GetEvaluationByID retrieves a single evaluation by its ID.
func (r *EvaluationRepository) GetEvaluationByID(id uint) (*models.Evaluation, error) {
	var evaluation models.Evaluation
	err := r.DB.Preload("Deliverable").Preload("Evaluator").Preload("Student").Preload("Rubric").Preload("CriterionEvaluations").First(&evaluation, id).Error
	return &evaluation, err
}

// GetEvaluationsByStudentID retrieves all evaluations for a given student ID.
func (r *EvaluationRepository) GetEvaluationsByStudentID(studentID uint) ([]models.Evaluation, error) {
	var evaluations []models.Evaluation
	err := r.DB.Where("student_id = ?", studentID).Preload("Deliverable").Preload("Evaluator").Preload("Student").Preload("Rubric").Find(&evaluations).Error
	return evaluations, err
}