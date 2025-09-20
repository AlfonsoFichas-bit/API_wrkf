package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// TaskRepository handles database operations for tasks.
type TaskRepository struct {
	DB *gorm.DB
}

// NewTaskRepository creates a new instance of TaskRepository.
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// CreateTask adds a new task to the database.
func (r *TaskRepository) CreateTask(task *models.Task) error {
	return r.DB.Create(task).Error
}

// GetTaskByID retrieves a single task by its ID, preloading related data.
func (r *TaskRepository) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.DB.Preload("UserStory").Preload("CreatedBy").Preload("AssignedTo").First(&task, id).Error
	return &task, err
}

// GetTasksByUserStoryID retrieves all tasks for a given user story ID.
func (r *TaskRepository) GetTasksByUserStoryID(userStoryID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Where("user_story_id = ?", userStoryID).Preload("CreatedBy").Preload("AssignedTo").Find(&tasks).Error
	return tasks, err
}

// UpdateTask updates an existing task in the database.
func (r *TaskRepository) UpdateTask(task *models.Task) error {
	return r.DB.Save(task).Error
}

// DeleteTask removes a task from the database by its ID.
func (r *TaskRepository) DeleteTask(id uint) error {
	return r.DB.Delete(&models.Task{}, id).Error
}

// DeleteTasksByUserStoryIDs deletes all tasks associated with a list of user story IDs.
func (r *TaskRepository) DeleteTasksByUserStoryIDs(tx *gorm.DB, storyIDs []uint) error {
	return tx.Where("user_story_id IN ?", storyIDs).Delete(&models.Task{}).Error
}
