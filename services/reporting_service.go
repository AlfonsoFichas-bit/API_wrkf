package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// ReportingService defines the interface for reporting-related business logic.
type ReportingService interface {
	CalculateProjectVelocity(projectID uint) (*models.VelocityReport, error)
	CalculateSprintBurndown(sprintID uint) (*models.BurndownReport, error)
	CalculateSprintCommitment(sprintID uint) (*models.CommitmentReport, error)
}

type reportingService struct {
	repo          storage.ReportingRepository
	userStoryRepo *storage.UserStoryRepository
	sprintRepo    *storage.SprintRepository
}

// NewReportingService creates a new instance of ReportingService.
func NewReportingService(repo storage.ReportingRepository, userStoryRepo *storage.UserStoryRepository, sprintRepo *storage.SprintRepository) ReportingService {
	return &reportingService{
		repo:          repo,
		userStoryRepo: userStoryRepo,
		sprintRepo:    sprintRepo,
	}
}

// CalculateProjectVelocity calculates the average velocity for a given project.
func (s *reportingService) CalculateProjectVelocity(projectID uint) (*models.VelocityReport, error) {
	sprints, err := s.repo.GetSprintsForVelocity(projectID)
	if err != nil {
		return nil, err
	}

	if len(sprints) == 0 {
		return &models.VelocityReport{
			ProjectID:         projectID,
			AverageVelocity:   0,
			SprintsConsidered: 0,
			VelocityPerSprint: []models.SprintVelocity{},
		}, nil
	}

	var totalPoints int
	var velocityPerSprint []models.SprintVelocity

	for _, sprint := range sprints {
		sprintPoints := 0
		for _, story := range sprint.UserStories {
			if story.Points != nil {
				sprintPoints += *story.Points
			}
		}
		totalPoints += sprintPoints
		velocityPerSprint = append(velocityPerSprint, models.SprintVelocity{
			SprintID:        sprint.ID,
			SprintName:      sprint.Name,
			CompletedPoints: sprintPoints,
		})
	}

	average := float64(totalPoints) / float64(len(sprints))

	report := &models.VelocityReport{
		ProjectID:         projectID,
		AverageVelocity:   average,
		SprintsConsidered: len(sprints),
		VelocityPerSprint: velocityPerSprint,
	}

	return report, nil
}

// CalculateSprintBurndown calculates the data points for a sprint's burndown chart.
func (s *reportingService) CalculateSprintBurndown(sprintID uint) (*models.BurndownReport, error) {
	sprint, err := s.repo.GetSprintForBurndown(sprintID)
	if err != nil {
		return nil, err
	}

	if sprint.StartDate == nil || sprint.EndDate == nil {
		return nil, errors.New("sprint must have a start and end date")
	}

	totalPoints := 0
	storyPointsMap := make(map[uint]int) // storyID -> points
	var storyIDs []uint
	for _, story := range sprint.UserStories {
		if story.Points != nil {
			totalPoints += *story.Points
			storyPointsMap[story.ID] = *story.Points
			storyIDs = append(storyIDs, story.ID)
		}
	}

	tasks, err := s.repo.GetTasksForUserStories(storyIDs)
	if err != nil {
		return nil, err
	}
	tasksByStoryID := make(map[uint][]models.Task)
	for _, task := range tasks {
		tasksByStoryID[task.UserStoryID] = append(tasksByStoryID[task.UserStoryID], task)
	}

	sprintDurationDays := int(sprint.EndDate.Sub(*sprint.StartDate).Hours()/24) + 1
	dailyIdealBurndown := float64(totalPoints) / float64(sprintDurationDays-1)

	// storyCompletionDate maps storyID to its completion date
	storyCompletionDate := make(map[uint]time.Time)
	for _, story := range sprint.UserStories {
		storyTasks := tasksByStoryID[story.ID]
		isCompleted := true
		var lastCompletion time.Time
		if len(storyTasks) == 0 {
			isCompleted = false // Story with no tasks is not considered done
		}
		for _, task := range storyTasks {
			if task.Status != "done" {
				isCompleted = false
				break
			}
			// Find the timestamp when the task was marked as 'done'
			for _, history := range task.History {
				if history.FieldName == "status" && history.NewValue == "done" && history.ChangedAt.After(lastCompletion) {
					lastCompletion = history.ChangedAt
				}
			}
		}
		if isCompleted {
			storyCompletionDate[story.ID] = lastCompletion
		}
	}

	burndownData := []models.BurndownPoint{}
	remainingPoints := float64(totalPoints)

	for i := 0; i < sprintDurationDays; i++ {
		currentDate := sprint.StartDate.AddDate(0, 0, i)

		// Calculate points completed on this day
		for storyID, completionTime := range storyCompletionDate {
			if isSameDay(completionTime, currentDate) {
				remainingPoints -= float64(storyPointsMap[storyID])
			}
		}

		// Ensure remaining points don't go below zero
		if remainingPoints < 0 {
			remainingPoints = 0
		}

		idealPoints := float64(totalPoints) - (float64(i) * dailyIdealBurndown)
		if idealPoints < 0 {
			idealPoints = 0
		}

		burndownData = append(burndownData, models.BurndownPoint{
			Date:            currentDate.Format("2006-01-02"),
			RemainingPoints: remainingPoints,
			IdealPoints:     idealPoints,
		})
	}

	report := &models.BurndownReport{
		SprintID:     sprint.ID,
		SprintName:   sprint.Name,
		TotalPoints:  totalPoints,
		BurndownData: burndownData,
	}

	return report, nil
}

func (s *reportingService) CalculateSprintCommitment(sprintID uint) (*models.CommitmentReport, error) {
	sprint, err := s.sprintRepo.GetSprintByID(sprintID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve sprint: %w", err)
	}

	stories, err := s.userStoryRepo.GetUserStoriesBySprintID(sprintID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve user stories for sprint: %w", err)
	}

	committedPoints := 0
	completedPoints := 0

	for _, story := range stories {
		if story.Points != nil {
			committedPoints += *story.Points
			if story.Status == "done" { // Assuming 'done' is the status for completed stories
				completedPoints += *story.Points
			}
		}
	}

	var completionRate float64
	if committedPoints > 0 {
		completionRate = (float64(completedPoints) / float64(committedPoints)) * 100
	}

	report := &models.CommitmentReport{
		SprintID:        sprint.ID,
		SprintName:      sprint.Name,
		CommittedPoints: committedPoints,
		CompletedPoints: completedPoints,
		CompletionRate:  completionRate,
	}

	return report, nil
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
