package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// EvaluationService handles business logic for evaluations.
type EvaluationService struct {
	evaluationRepo *storage.EvaluationRepository
}

// NewEvaluationService creates a new instance of EvaluationService.
func NewEvaluationService(evaluationRepo *storage.EvaluationRepository) *EvaluationService {
	return &EvaluationService{evaluationRepo: evaluationRepo}
}

// CreateEvaluation creates a new evaluation.
func (s *EvaluationService) CreateEvaluation(evaluation *models.Evaluation) error {
	return s.evaluationRepo.CreateEvaluation(evaluation)
}

// GetEvaluationByID retrieves a single evaluation by its ID.
func (s *EvaluationService) GetEvaluationByID(id uint) (*models.Evaluation, error) {
	return s.evaluationRepo.GetEvaluationByID(id)
}

// GetEvaluationsByStudentID retrieves all evaluations for a given student ID.
func (s *EvaluationService) GetEvaluationsByStudentID(studentID uint) ([]models.Evaluation, error) {
	return s.evaluationRepo.GetEvaluationsByStudentID(studentID)
}