package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// ReportingRepository defines the interface for reporting-related database operations.
type ReportingRepository interface {
	GetSprintsForVelocity(projectID uint) ([]models.Sprint, error)
	GetSprintForBurndown(sprintID uint) (*models.Sprint, error)
	GetTasksForUserStories(storyIDs []uint) ([]models.Task, error)
}

type reportingRepository struct {
	db *gorm.DB
}

// NewReportingRepository creates a new instance of ReportingRepository.
func NewReportingRepository(db *gorm.DB) ReportingRepository {
	return &reportingRepository{db: db}
}

// GetSprintsForVelocity fetches completed sprints for a project to calculate velocity.
// It preloads the user stories to access their points.
func (r *reportingRepository) GetSprintsForVelocity(projectID uint) ([]models.Sprint, error) {
	var sprints []models.Sprint
	err := r.db.
		Preload("UserStories", "status = ?", "done").
		Where("project_id = ? AND status = ?", projectID, "completed").
		Order("end_date asc").
		Find(&sprints).Error
	return sprints, err
}

// GetSprintForBurndown fetches a sprint with its user stories.
// The tasks are fetched in a separate call.
func (r *reportingRepository) GetSprintForBurndown(sprintID uint) (*models.Sprint, error) {
	var sprint models.Sprint
	err := r.db.
		Preload("UserStories").
		First(&sprint, sprintID).Error
	return &sprint, err
}

// GetTasksForUserStories fetches all tasks and their history for a given list of user story IDs.
func (r *reportingRepository) GetTasksForUserStories(storyIDs []uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.
		Preload("History").
		Where("user_story_id IN ?", storyIDs).
		Find(&tasks).Error
	return tasks, err
}
