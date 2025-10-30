package services

import (
	"fmt"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// BurndownDataPoint represents a single point in the burndown chart.
type BurndownDataPoint struct {
	Date         string  `json:"date"`
	IdealPoints  float64 `json:"idealPoints"`
	ActualPoints float64 `json:"actualPoints"`
}

// BurndownChart represents the data needed to render a burndown chart.
type BurndownChart struct {
	SprintName string              `json:"sprintName"`
	DataPoints []BurndownDataPoint `json:"dataPoints"`
}

// BurndownService is responsible for calculating burndown chart data.
type BurndownService struct {
	repo storage.IBurndownRepository
}

// NewBurndownService creates a new instance of BurndownService.
func NewBurndownService(repo storage.IBurndownRepository) *BurndownService {
	return &BurndownService{repo: repo}
}

// GenerateBurndownChart generates the data for a sprint's burndown chart.
func (s *BurndownService) GenerateBurndownChart(sprintID uint) (*BurndownChart, error) {
	sprint, err := s.repo.GetSprintWithTasks(sprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprint: %w", err)
	}

	if sprint.StartDate == nil || sprint.EndDate == nil {
		return nil, fmt.Errorf("sprint must have a start and end date")
	}

	histories, err := s.repo.GetTaskHistoryForSprint(sprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task history: %w", err)
	}

	totalStoryPoints := 0
	for _, us := range sprint.UserStories {
		for _, task := range us.Tasks {
			totalStoryPoints += task.StoryPoints
		}
	}

	sprintDurationDays := int(sprint.EndDate.Sub(*sprint.StartDate).Hours()/24) + 1
	if sprintDurationDays <= 0 {
		sprintDurationDays = 1
	}

	// Avoid division by zero if sprint duration is 1 day
	idealPointsPerDay := 0.0
	if sprintDurationDays > 1 {
		idealPointsPerDay = float64(totalStoryPoints) / float64(sprintDurationDays-1)
	}

	pointsPerDay := make(map[string]int)
	for _, history := range histories {
		if history.FieldName == "status" && history.NewValue == string(models.StatusDone) {
			// Find the task to get its story points
			for _, us := range sprint.UserStories {
				for _, task := range us.Tasks {
					if task.ID == history.TaskID {
						day := history.ChangedAt.Format("2006-01-02")
						pointsPerDay[day] += task.StoryPoints
						break
					}
				}
			}
		}
	}

	var dataPoints []BurndownDataPoint
	remainingPoints := float64(totalStoryPoints)
	currentDate := *sprint.StartDate
	endDate := *sprint.EndDate

	for i := 0; i < sprintDurationDays; i++ {
		dateStr := currentDate.Format("2006-01-02")

		// Points completed on day 0 are not subtracted from the total
		if i > 0 {
			dayBefore := currentDate.AddDate(0, 0, -1).Format("2006-01-02")
			if pointsCompleted, ok := pointsPerDay[dayBefore]; ok {
				remainingPoints -= float64(pointsCompleted)
			}
		}


		ideal := float64(totalStoryPoints) - (float64(i) * idealPointsPerDay)
		if ideal < 0 {
			ideal = 0
		}

		// On day 0, the actual points should be the total points
		if i == 0 {
			dataPoints = append(dataPoints, BurndownDataPoint{
				Date:         dateStr,
				IdealPoints:  float64(totalStoryPoints),
				ActualPoints: float64(totalStoryPoints),
			})
		} else {
			dataPoints = append(dataPoints, BurndownDataPoint{
				Date:         dateStr,
				IdealPoints:  ideal,
				ActualPoints: remainingPoints,
			})
		}


		currentDate = currentDate.AddDate(0, 0, 1)
		if currentDate.After(endDate.Add(time.Hour * 24)) {
			break
		}
	}

	return &BurndownChart{
		SprintName: sprint.Name,
		DataPoints: dataPoints,
	}, nil
}
