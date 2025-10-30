package mocks

import (
	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/mock"
)

type MockRubricRepository struct {
	mock.Mock
}

func (m *MockRubricRepository) CreateRubric(rubric *models.Rubric) error {
	args := m.Called(rubric)
	return args.Error(0)
}

func (m *MockRubricRepository) GetAllRubrics(filters map[string]interface{}) ([]models.Rubric, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Rubric), args.Error(1)
}

func (m *MockRubricRepository) GetRubricByID(id uint) (*models.Rubric, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rubric), args.Error(1)
}

func (m *MockRubricRepository) UpdateRubric(rubric *models.Rubric) error {
	args := m.Called(rubric)
	return args.Error(0)
}

func (m *MockRubricRepository) DeleteRubric(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
