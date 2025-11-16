package services

import (
	"fmt"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// SprintService handles the business logic for sprints.
type SprintService struct {
	Repo             *storage.SprintRepository
	ReportingService ReportingService
	ActivityService  *ActivityService
}

// NewSprintService creates a new instance of SprintService.
func NewSprintService(repo *storage.SprintRepository, reportingService ReportingService, activityService *ActivityService) *SprintService {
	return &SprintService{Repo: repo, ReportingService: reportingService, ActivityService: activityService}
}

// CreateSprint handles the business logic for creating a new sprint.
func (s *SprintService) CreateSprint(sprint *models.Sprint, projectID uint, creatorID uint) error {
	sprint.ProjectID = projectID
	sprint.CreatedByID = creatorID
	return s.Repo.CreateSprint(sprint)
}

// SprintWithStats contains a sprint and its calculated statistics.
type SprintWithStats struct {
	models.Sprint
	TotalTasks      int     `json:"totalTasks"`
	CompletedTasks  int     `json:"completedTasks"`
	InProgressTasks int     `json:"inProgressTasks"`
	PendingTasks    int     `json:"pendingTasks"`
	Progress        float64 `json:"progress"`
	Velocity        int     `json:"velocity"`
	// TeamVelocity will be handled in the handler/response layer
}

// GetSprintsByProjectID retrieves all sprints for a specific project with calculated stats.
func (s *SprintService) GetSprintsByProjectID(projectID uint) ([]SprintWithStats, *models.Sprint, error) {
	sprints, err := s.Repo.GetSprintsByProjectID(projectID)
	if err != nil {
		return nil, nil, err
	}

	if len(sprints) == 0 {
		return []SprintWithStats{}, nil, nil
	}

	sprintIDs := make([]uint, len(sprints))
	for i, sprint := range sprints {
		sprintIDs[i] = sprint.ID
	}

	taskCounts, err := s.Repo.GetTaskCountsForSprints(sprintIDs)
	if err != nil {
		return nil, nil, err
	}

	sprintsWithStats := make([]SprintWithStats, len(sprints))
	var activeSprint *models.Sprint
	for i, sprint := range sprints {
		stats := taskCounts[sprint.ID]
		completed := stats[models.StatusDone]
		inProgress := stats[models.StatusInProgress]
		inReview := stats[models.StatusInReview]
		todo := stats[models.StatusTodo]

		total := completed + inProgress + inReview + todo
		var progress float64
		if total > 0 {
			progress = (float64(completed) / float64(total)) * 100
		}

		// Calculate velocity (completed story points) for the sprint
		commitmentReport, err := s.ReportingService.CalculateSprintCommitment(sprint.ID)
		var velocity int
		if err != nil {
			// If there's an error (e.g., sprint not found), velocity is 0
			velocity = 0
		} else {
			velocity = commitmentReport.CompletedPoints
		}


		sprintsWithStats[i] = SprintWithStats{
			Sprint:          sprint,
			TotalTasks:      total,
			CompletedTasks:  completed,
			InProgressTasks: inProgress + inReview, // Combining in_progress and in_review
			PendingTasks:    todo,
			Progress:        progress,
			Velocity:        velocity,
		}
		if sprint.Status == "active" {
			activeSprint = &sprints[i]
		}
	}

	return sprintsWithStats, activeSprint, nil
}

// GetSprintByID retrieves a single sprint by its ID.
func (s *SprintService) GetSprintByID(id uint) (*models.Sprint, error) {
	return s.Repo.GetSprintByID(id)
}

// UpdateSprint handles the business logic for updating a sprint.
func (s *SprintService) UpdateSprint(sprint *models.Sprint) error {
	// In the future, you could add permission checks here.
	return s.Repo.UpdateSprint(sprint)
}

// DeleteSprint handles the business logic for deleting a sprint.
func (s *SprintService) DeleteSprint(id uint) error {
	// In future, you could add logic here to move user stories back to the backlog.
	return s.Repo.DeleteSprint(id)
}

// GetSprintTasks retrieves all tasks for a specific sprint with their relationships.
func (s *SprintService) GetSprintTasks(sprintID uint) ([]models.Task, error) {
	return s.Repo.GetSprintTasks(sprintID)
}

// UpdateSprintStatus updates the status of a sprint.
func (s *SprintService) UpdateSprintStatus(sprintID uint, status string) error {
	if status == "active" {
		sprint, err := s.Repo.GetSprintByID(sprintID)
		if err != nil {
			return fmt.Errorf("sprint not found")
		}
		activeSprint, err := s.Repo.GetActiveSprint(sprint.ProjectID)
		if err == nil && activeSprint != nil && activeSprint.ID != sprintID {
			return fmt.Errorf("another sprint is already active in this project")
		}

		// --- Create Activity ---
		description := fmt.Sprintf("ha iniciado el sprint '%s'", sprint.Name)
		// We don't have the user ID who started the sprint. For now, we'll use the sprint creator's ID.
		s.ActivityService.CreateActivity(
			"sprint_started",
			"sprint",
			sprintID,
			sprint.CreatedByID,
			sprint.ProjectID,
			description,
		)
		// --- End Activity ---
	}
	return s.Repo.UpdateSprintStatus(sprintID, status)
}
