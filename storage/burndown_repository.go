package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// IBurndownRepository defines the interface for accessing the data needed
// to generate a burndown chart.
type IBurndownRepository interface {
	GetSprintWithTasks(sprintID uint) (*models.Sprint, error)
	GetTaskHistoryForSprint(sprintID uint) ([]models.TaskHistory, error)
}

// BurndownRepository is the concrete implementation of IBurndownRepository.
type BurndownRepository struct {
	db *gorm.DB
}

// NewBurndownRepository creates a new instance of BurndownRepository.
func NewBurndownRepository(db *gorm.DB) *BurndownRepository {
	return &BurndownRepository{db: db}
}

// GetSprintWithTasks retrieves a sprint by its ID, preloading its user stories and their tasks.
func (r *BurndownRepository) GetSprintWithTasks(sprintID uint) (*models.Sprint, error) {
	var sprint models.Sprint
	// We need to preload the full hierarchy to be able to calculate total story points.
	err := r.db.Preload("UserStories.Tasks").First(&sprint, sprintID).Error
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}

// GetTaskHistoryForSprint retrieves all task history records for a given sprint.
func (r *BurndownRepository) GetTaskHistoryForSprint(sprintID uint) ([]models.TaskHistory, error) {
	var histories []models.TaskHistory
	// We join across the tables to find all task histories associated with the sprint.
	err := r.db.
		Joins("JOIN tasks ON tasks.id = task_histories.task_id").
		Joins("JOIN user_stories ON user_stories.id = tasks.user_story_id").
		Where("user_stories.sprint_id = ?", sprintID).
		Order("task_histories.changed_at asc").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}
	return histories, nil
}
