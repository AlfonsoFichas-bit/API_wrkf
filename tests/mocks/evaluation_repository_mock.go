package mocks

import (
	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/mock"
)

type MockEvaluationRepository struct {
	mock.Mock
}

func (m *MockEvaluationRepository) CreateEvaluation(evaluation *models.Evaluation) error {
	args := m.Called(evaluation)
	return args.Error(0)
}

func (m *MockEvaluationRepository) GetEvaluationByID(id uint) (*models.Evaluation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Evaluation), args.Error(1)
}

func (m *MockEvaluationRepository) GetEvaluationsByStudent(studentID uint) ([]models.Evaluation, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Evaluation), args.Error(1)
}

func (m *MockEvaluationRepository) UpdateEvaluation(evaluation *models.Evaluation) error {
	args := m.Called(evaluation)
	return args.Error(0)
}
