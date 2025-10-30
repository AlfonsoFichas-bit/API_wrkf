package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

type IEvaluationRepository interface {
	CreateEvaluation(evaluation *models.Evaluation) error
	GetEvaluationByID(id uint) (*models.Evaluation, error)
	GetEvaluationsByStudent(studentID uint) ([]models.Evaluation, error)
	UpdateEvaluation(evaluation *models.Evaluation) error
}

type EvaluationRepository struct {
	db *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) *EvaluationRepository {
	return &EvaluationRepository{db: db}
}

func (r *EvaluationRepository) CreateEvaluation(evaluation *models.Evaluation) error {
	return r.db.Create(evaluation).Error
}

func (r *EvaluationRepository) GetEvaluationByID(id uint) (*models.Evaluation, error) {
	var evaluation models.Evaluation
	err := r.db.Preload("Grades").First(&evaluation, id).Error
	return &evaluation, err
}

func (r *EvaluationRepository) GetEvaluationsByStudent(studentID uint) ([]models.Evaluation, error) {
	var evaluations []models.Evaluation
	err := r.db.Where("student_id = ?", studentID).Find(&evaluations).Error
	return evaluations, err
}

func (r *EvaluationRepository) UpdateEvaluation(evaluation *models.Evaluation) error {
	return r.db.Save(evaluation).Error
}
