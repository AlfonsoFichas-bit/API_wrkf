package services_test

import (
	"testing"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBurndownRepository is a mock implementation of IBurndownRepository for testing.
type MockBurndownRepository struct {
	mock.Mock
}

func (m *MockBurndownRepository) GetSprintWithTasks(sprintID uint) (*models.Sprint, error) {
	args := m.Called(sprintID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Sprint), args.Error(1)
}

func (m *MockBurndownRepository) GetTaskHistoryForSprint(sprintID uint) ([]models.TaskHistory, error) {
	args := m.Called(sprintID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TaskHistory), args.Error(1)
}

func TestGenerateBurndownChart_Success(t *testing.T) {
	mockRepo := new(MockBurndownRepository)
	burndownService := services.NewBurndownService(mockRepo)

	// --- Test Data ---
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC) // 5-day sprint

	sprint := &models.Sprint{
		ID:        1,
		Name:      "Test Sprint",
		StartDate: &startDate,
		EndDate:   &endDate,
		UserStories: []models.UserStory{
			{
				ID: 1,
				Tasks: []models.Task{
					{ID: 1, StoryPoints: 5, UserStoryID: 1},
					{ID: 2, StoryPoints: 3, UserStoryID: 1},
					{ID: 3, StoryPoints: 2, UserStoryID: 1},
				},
			},
		},
	}

	histories := []models.TaskHistory{
		// Day 1 (end of day): Task 1 (5 points) is completed
		{TaskID: 1, FieldName: "status", NewValue: "done", ChangedAt: startDate},
		// Day 3 (end of day): Task 2 (3 points) is completed
		{TaskID: 2, FieldName: "status", NewValue: "done", ChangedAt: startDate.AddDate(0, 0, 2)},
	}

	// --- Mock Expectations ---
	mockRepo.On("GetSprintWithTasks", uint(1)).Return(sprint, nil)
	mockRepo.On("GetTaskHistoryForSprint", uint(1)).Return(histories, nil)

	// --- Execute ---
	chart, err := burndownService.GenerateBurndownChart(1)

	// --- Assertions ---
	assert.NoError(t, err)
	assert.NotNil(t, chart)
	assert.Len(t, chart.DataPoints, 5) // 5 days in the sprint

	// Expected values: Total 10 points, 5 days. Ideal burn is 2.5 points/day.
	// Day 1 (Jan 1): Ideal 10, Actual 10.
	assert.Equal(t, "2023-01-01", chart.DataPoints[0].Date)
	assert.InDelta(t, 10, chart.DataPoints[0].IdealPoints, 0.01)
	assert.InDelta(t, 10, chart.DataPoints[0].ActualPoints, 0.01, "Day 1 actual should be total points")

	// Day 2 (Jan 2): Ideal 7.5, Actual 5. (5 points were burned on Day 1)
	assert.Equal(t, "2023-01-02", chart.DataPoints[1].Date)
	assert.InDelta(t, 7.5, chart.DataPoints[1].IdealPoints, 0.01)
	assert.InDelta(t, 5, chart.DataPoints[1].ActualPoints, 0.01, "Day 2 actual should reflect points burned on Day 1")

	// Day 3 (Jan 3): Ideal 5.0, Actual 5.
	assert.Equal(t, "2023-01-03", chart.DataPoints[2].Date)
	assert.InDelta(t, 5.0, chart.DataPoints[2].IdealPoints, 0.01)
	assert.InDelta(t, 5, chart.DataPoints[2].ActualPoints, 0.01, "Day 3 actual should be unchanged")

	// Day 4 (Jan 4): Ideal 2.5, Actual 2. (3 points were burned on Day 3)
	assert.Equal(t, "2023-01-04", chart.DataPoints[3].Date)
	assert.InDelta(t, 2.5, chart.DataPoints[3].IdealPoints, 0.01)
	assert.InDelta(t, 2, chart.DataPoints[3].ActualPoints, 0.01, "Day 4 actual should reflect points burned on Day 3")

	// Day 5 (Jan 5): Ideal 0, Actual 2.
	assert.Equal(t, "2023-01-05", chart.DataPoints[4].Date)
	assert.InDelta(t, 0, chart.DataPoints[4].IdealPoints, 0.01)
	assert.InDelta(t, 2, chart.DataPoints[4].ActualPoints, 0.01, "Day 5 actual should be unchanged")

	mockRepo.AssertExpectations(t)
}
