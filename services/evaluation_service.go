package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

type IEvaluationService interface {
	CreateEvaluation(evaluation *models.Evaluation) error
	GetEvaluationByID(id uint) (*models.Evaluation, error)
	GetEvaluationsByStudent(studentID uint) ([]models.Evaluation, error)
	CalculateFinalGrade(evaluationID uint) (float64, error)
}

type EvaluationService struct {
	evaluationRepo storage.IEvaluationRepository
	rubricRepo     storage.IRubricRepository
}

func NewEvaluationService(evaluationRepo storage.IEvaluationRepository, rubricRepo storage.IRubricRepository) *EvaluationService {
	return &EvaluationService{
		evaluationRepo: evaluationRepo,
		rubricRepo:     rubricRepo,
	}
}

func (s *EvaluationService) CreateEvaluation(evaluation *models.Evaluation) error {
	return s.evaluationRepo.CreateEvaluation(evaluation)
}

func (s *EvaluationService) GetEvaluationByID(id uint) (*models.Evaluation, error) {
	return s.evaluationRepo.GetEvaluationByID(id)
}

func (s *EvaluationService) GetEvaluationsByStudent(studentID uint) ([]models.Evaluation, error) {
	return s.evaluationRepo.GetEvaluationsByStudent(studentID)
}

func (s *EvaluationService) CalculateFinalGrade(evaluationID uint) (float64, error) {
	evaluation, err := s.evaluationRepo.GetEvaluationByID(evaluationID)
	if err != nil {
		return 0, err
	}

	rubric, err := s.rubricRepo.GetRubricByID(evaluation.RubricID)
	if err != nil {
		return 0, err
	}

	var totalScore uint
	var maxPoints uint
	for _, criterion := range rubric.Criteria {
		var criterionMaxPoints uint
		for _, level := range criterion.Levels {
			if level.Points > criterionMaxPoints {
				criterionMaxPoints = level.Points
			}
		}
		maxPoints += criterionMaxPoints

		for _, grade := range evaluation.Grades {
			if grade.CriterionID == criterion.ID {
				totalScore += grade.Score
				break
			}
		}
	}

	if maxPoints == 0 {
		return 0, nil
	}

	finalGrade := (float64(totalScore) / float64(maxPoints)) * 100
	evaluation.FinalGrade = finalGrade

	if err := s.evaluationRepo.UpdateEvaluation(evaluation); err != nil {
        return 0, err
    }

	return finalGrade, nil
}
