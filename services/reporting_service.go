package services

import (
	"encoding/json"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// ReportingService handles business logic for reports.
type ReportingService struct {
	metricService *MetricService
	sprintRepo    *storage.SprintRepository
}

// NewReportingService creates a new instance of ReportingService.
func NewReportingService(metricService *MetricService, sprintRepo *storage.SprintRepository) *ReportingService {
	return &ReportingService{metricService: metricService, sprintRepo: sprintRepo}
}

// GenerateProjectReport generates a report for a given project.
func (s *ReportingService) GenerateProjectReport(projectID uint) (*models.Report, error) {
	sprints, err := s.sprintRepo.GetSprintsByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	burndownCharts := make(map[string][]BurndownChartPoint)
	for _, sprint := range sprints {
		burndownChart, err := s.metricService.GetBurndownChart(sprint.ID)
		if err != nil {
			return nil, err
		}
		burndownCharts[sprint.Name] = burndownChart
	}

	teamVelocity, err := s.metricService.GetTeamVelocity(projectID)
	if err != nil {
		return nil, err
	}

	report := &models.Report{
		Title:       "Project Report",
		Description: "A report of the project's progress.",
		Type:        "project",
		ProjectID:   &projectID,
	}

	data := make(map[string]interface{})
	data["burndownCharts"] = burndownCharts
	data["teamVelocity"] = teamVelocity

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	report.Data = jsonData

	return report, nil
}