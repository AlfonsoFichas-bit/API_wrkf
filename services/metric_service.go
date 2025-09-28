package services

import (
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// MetricService handles business logic for metrics.
type MetricService struct {
	sprintRepo *storage.SprintRepository
	taskRepo   *storage.TaskRepository
}

// NewMetricService creates a new instance of MetricService.
func NewMetricService(sprintRepo *storage.SprintRepository, taskRepo *storage.TaskRepository) *MetricService {
	return &MetricService{sprintRepo: sprintRepo, taskRepo: taskRepo}
}

// BurndownChartPoint represents a single point in the burndown chart.
type BurndownChartPoint struct {
	Date          time.Time `json:"date"`
	RemainingWork float32   `json:"remainingWork"`
}

// GetBurndownChart calculates the burndown chart for a given sprint.
func (s *MetricService) GetBurndownChart(sprintID uint) ([]BurndownChartPoint, error) {
	sprint, err := s.sprintRepo.GetSprintByID(sprintID)
	if err != nil {
		return nil, err
	}

	tasks, err := s.taskRepo.GetTasksBySprintID(sprintID)
	if err != nil {
		return nil, err
	}

	burndownChart := []BurndownChartPoint{}
	if sprint.StartDate == nil || sprint.EndDate == nil {
		return burndownChart, nil
	}

	totalWork := float32(0)
	for _, task := range tasks {
		if task.EstimatedHours != nil {
			totalWork += *task.EstimatedHours
		}
	}

	for d := *sprint.StartDate; d.Before(*sprint.EndDate); d = d.AddDate(0, 0, 1) {
		remainingWork := totalWork
		for _, task := range tasks {
			if task.Status == string(models.StatusDone) && task.UpdatedAt.Before(d) {
				if task.EstimatedHours != nil {
					remainingWork -= *task.EstimatedHours
				}
			}
		}
		burndownChart = append(burndownChart, BurndownChartPoint{Date: d, RemainingWork: remainingWork})
	}

	return burndownChart, nil
}

// TeamVelocity represents the team velocity for a given project.
type TeamVelocity struct {
	SprintName string `json:"sprintName"`
	CompletedWork int    `json:"completedWork"`
}

// GetTeamVelocity calculates the team velocity for a given project.
func (s *MetricService) GetTeamVelocity(projectID uint) ([]TeamVelocity, error) {
	sprints, err := s.sprintRepo.GetSprintsByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	teamVelocity := []TeamVelocity{}
	for _, sprint := range sprints {
		if sprint.Status == "completed" {
			tasks, err := s.taskRepo.GetTasksBySprintID(sprint.ID)
			if err != nil {
				return nil, err
			}

			completedWork := 0
			for _, task := range tasks {
				if task.Status == string(models.StatusDone) {
					completedWork++
				}
			}
			teamVelocity = append(teamVelocity, TeamVelocity{SprintName: sprint.Name, CompletedWork: completedWork})
		}
	}

	return teamVelocity, nil
}

// WorkDistribution represents the work distribution for a given sprint.
type WorkDistribution struct {
	UserName string `json:"userName"`
	TaskCount int    `json:"taskCount"`
}

// GetWorkDistribution calculates the work distribution for a given sprint.
func (s *MetricService) GetWorkDistribution(sprintID uint) ([]WorkDistribution, error) {
	tasks, err := s.taskRepo.GetTasksBySprintID(sprintID)
	if err != nil {
		return nil, err
	}

	workDistribution := make(map[string]int)
	for _, task := range tasks {
		if task.AssignedTo != nil {
			workDistribution[task.AssignedTo.Nombre]++
		}
	}

	result := []WorkDistribution{}
	for userName, taskCount := range workDistribution {
		result = append(result, WorkDistribution{UserName: userName, TaskCount: taskCount})
	}

	return result, nil
}