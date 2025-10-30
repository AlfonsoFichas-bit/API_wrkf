package services

import (
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEvaluationService_GetEvaluationByID(t *testing.T) {
	mockEvaluationRepo := new(mocks.MockEvaluationRepository)
	mockRubricRepo := new(mocks.MockRubricRepository)
	service := services.NewEvaluationService(mockEvaluationRepo, mockRubricRepo)

	evaluation := &models.Evaluation{StudentID: 1, ProjectID: 1, TaskID: 1}

	mockEvaluationRepo.On("GetEvaluationByID", uint(1)).Return(evaluation, nil)

	result, err := service.GetEvaluationByID(uint(1))

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.StudentID)
	mockEvaluationRepo.AssertExpectations(t)
}

func TestEvaluationService_CreateEvaluation(t *testing.T) {
	mockEvaluationRepo := new(mocks.MockEvaluationRepository)
	mockRubricRepo := new(mocks.MockRubricRepository)
	service := services.NewEvaluationService(mockEvaluationRepo, mockRubricRepo)

	evaluation := &models.Evaluation{StudentID: 1, ProjectID: 1, TaskID: 1}

	mockEvaluationRepo.On("CreateEvaluation", evaluation).Return(nil)

	err := service.CreateEvaluation(evaluation)

	assert.NoError(t, err)
	mockEvaluationRepo.AssertExpectations(t)
}
