package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// SprintRepository handles database operations for sprints.
type SprintRepository struct {
	DB *gorm.DB
}

// NewSprintRepository creates a new instance of SprintRepository.
func NewSprintRepository(db *gorm.DB) *SprintRepository {
	return &SprintRepository{DB: db}
}

// CreateSprint adds a new sprint to the database.
func (r *SprintRepository) CreateSprint(sprint *models.Sprint) error {
	return r.DB.Create(sprint).Error
}

// GetSprintsByProjectID retrieves all sprints for a given project ID.
func (r *SprintRepository) GetSprintsByProjectID(projectID uint) ([]models.Sprint, error) {
	var sprints []models.Sprint
	err := r.DB.Where("project_id = ?", projectID).Preload("CreatedBy").Find(&sprints).Error
	return sprints, err
}

// GetSprintByID retrieves a single sprint by its ID.
func (r *SprintRepository) GetSprintByID(id uint) (*models.Sprint, error) {
	var sprint models.Sprint
	err := r.DB.Preload("CreatedBy").Preload("Project").First(&sprint, id).Error
	return &sprint, err
}

// UpdateSprint updates an existing sprint in the database.
func (r *SprintRepository) UpdateSprint(sprint *models.Sprint) error {
	return r.DB.Save(sprint).Error
}

// DeleteSprint removes a sprint from the database by its ID.
func (r *SprintRepository) DeleteSprint(id uint) error {
	result := r.DB.Delete(&models.Sprint{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteSprintsByProjectID deletes all sprints associated with a project.
func (r *SprintRepository) DeleteSprintsByProjectID(tx *gorm.DB, projectID uint) error {
	return tx.Where("project_id = ?", projectID).Delete(&models.Sprint{}).Error
}

// GetSprintTasks retrieves all tasks for a specific sprint with their relationships.
func (r *SprintRepository) GetSprintTasks(sprintID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.
		Joins("JOIN user_stories ON tasks.user_story_id = user_stories.id").
		Where("user_stories.sprint_id = ?", sprintID).
		Preload("AssignedTo").
		Preload("CreatedBy").
		Preload("UserStory").
		Preload("UserStory.Project").
		Order("tasks.created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// UpdateSprintStatus updates the status of a sprint.
func (r *SprintRepository) UpdateSprintStatus(sprintID uint, status string) error {
	validStatuses := []string{"planned", "active", "completed", "closed"}
	isValidStatus := false
	for _, s := range validStatuses {
		if s == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return gorm.ErrInvalidData // Or a custom error type
	}

	result := r.DB.Model(&models.Sprint{}).Where("id = ?", sprintID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetActiveSprint retrieves the currently active sprint for a project.
func (r *SprintRepository) GetActiveSprint(projectID uint) (*models.Sprint, error) {
	var sprint models.Sprint
	err := r.DB.
		Where("project_id = ? AND status = ?", projectID, "active").
		Preload("CreatedBy").
		Preload("Project").
		First(&sprint).Error
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}

// TaskStats holds the count of tasks for a given status.
type TaskStats struct {
	SprintID uint
	Status   models.TaskStatus
	Count    int
}

// GetTaskCountsForSprints retrieves the count of tasks grouped by status for a list of sprint IDs.
func (r *SprintRepository) GetTaskCountsForSprints(sprintIDs []uint) (map[uint]map[models.TaskStatus]int, error) {
	var stats []TaskStats
	err := r.DB.Table("tasks").
		Select("user_stories.sprint_id, tasks.status, COUNT(tasks.id) as count").
		Joins("JOIN user_stories ON tasks.user_story_id = user_stories.id").
		Where("user_stories.sprint_id IN ?", sprintIDs).
		Group("user_stories.sprint_id, tasks.status").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// Process the flat results into a nested map
	result := make(map[uint]map[models.TaskStatus]int)
	for _, stat := range stats {
		if _, ok := result[stat.SprintID]; !ok {
			result[stat.SprintID] = make(map[models.TaskStatus]int)
		}
		result[stat.SprintID][stat.Status] = stat.Count
	}

	return result, nil
}
