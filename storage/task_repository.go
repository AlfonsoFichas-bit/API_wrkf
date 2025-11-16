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

// UpdateTask is a robust method to save a task. It explicitly specifies which
// fields should be updated, preventing GORM from accidentally nullifying associations.
func (r *TaskRepository) UpdateTask(task *models.Task) error {
	// By using `Select`, we tell GORM exactly which fields we intend to update.
	// This is the definitive fix for both the Kanban and reassignment bugs.
	return r.DB.Select(
		"Title",
		"Description",
		"Status",
		"AssignedToID",
		"EstimatedHours",
		"SpentHours",
		"IsDeliverable",
	).Save(task).Error
}

// DeleteTask removes a task and its dependencies from the database by its ID.
func (r *TaskRepository) DeleteTask(id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// First, delete dependent records to avoid foreign key violations.
		if err := tx.Where("task_id = ?", id).Delete(&models.TaskHistory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("task_id = ?", id).Delete(&models.TaskComment{}).Error; err != nil {
			return err
		}

		// Then, delete the task itself.
		if err := tx.Delete(&models.Task{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteTasksByUserStoryIDs deletes all tasks associated with a list of user story IDs.
func (r *TaskRepository) DeleteTasksByUserStoryIDs(tx *gorm.DB, storyIDs []uint) error {
	return tx.Where("user_story_id IN ?", storyIDs).Delete(&models.Task{}).Error
}

// AddComment adds a new comment to the database.
func (r *TaskRepository) AddComment(comment *models.TaskComment) error {
	return r.DB.Create(comment).Error
}

// GetProjectIDForTask finds the ProjectID for a given task by traversing up.
func (r *TaskRepository) GetProjectIDForTask(taskID uint) (uint, error) {
	var task models.Task
	// Select only the user_story_id field for efficiency
	if err := r.DB.Select("user_story_id").First(&task, taskID).Error; err != nil {
		return 0, err
	}

	var userStory models.UserStory
	// Select only the project_id field for efficiency
	if err := r.DB.Select("project_id").First(&userStory, task.UserStoryID).Error; err != nil {
		return 0, err
	}

	return userStory.ProjectID, nil
}

// GetCommentsByTaskID retrieves all comments for a given task, ordered by creation time.
func (r *TaskRepository) GetCommentsByTaskID(taskID uint) ([]models.TaskComment, error) {
	var comments []models.TaskComment
	err := r.DB.
		Where("task_id = ?", taskID).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

// UpdateTaskWithHistory updates a task's status and records the change in the history table.
func (r *TaskRepository) UpdateTaskWithHistory(task *models.Task, oldStatus, newStatus models.TaskStatus, changedByID uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Update only the status field. Using Model().Update() is safer here.
		if err := tx.Model(task).Update("status", newStatus).Error; err != nil {
			return err
		}

		// 2. Create the history record with the ID of the user who made the change.
		historyRecord := &models.TaskHistory{
			TaskID:      task.ID,
			ChangedByID: changedByID,
			FieldName:   "status",
			OldValue:    string(oldStatus),
			NewValue:    string(newStatus),
		}
		if err := tx.Create(historyRecord).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetTasksByUserID retrieves tasks for a user with optional filters.
func (r *TaskRepository) GetTasksByUserID(userID uint, projectID *uint, status *string, limit, offset int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.DB.Model(&models.Task{}).
		Joins("JOIN user_stories ON user_stories.id = tasks.user_story_id").
		Where("tasks.assigned_to_id = ?", userID)

	if projectID != nil {
		query = query.Where("user_stories.project_id = ?", *projectID)
	}

	if status != nil && *status != "" {
		query = query.Where("tasks.status = ?", *status)
	}

	// Count total records for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply limit, offset, and ordering
	err := query.
		Preload("UserStory").
		Preload("UserStory.Project").
		Order("tasks.due_date ASC").
		Order("CASE tasks.priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 ELSE 4 END").
		Limit(limit).
		Offset(offset).
		Find(&tasks).Error

	return tasks, total, err
}

// GetPendingEvaluations retrieves tasks submitted for evaluation for specific projects.
func (r *TaskRepository) GetPendingEvaluations(projectIDs []uint, limit int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.DB.Model(&models.Task{}).
		Joins("JOIN user_stories ON user_stories.id = tasks.user_story_id").
		Where("tasks.submitted_for_evaluation = ?", true).
		Where("user_stories.project_id IN ?", projectIDs)

	// Count total records for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("SubmittedBy").
		Preload("UserStory.Project").
		Order("tasks.submitted_at ASC").
		Limit(limit).
		Find(&tasks).Error

	return tasks, total, err
}
