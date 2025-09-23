package services

import (
	"fmt"
	"log"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// TaskService handles the business logic for tasks.
type TaskService struct {
	Repo                *storage.TaskRepository
	ProjectService      *ProjectService // Dependency to check project membership
	NotificationService *NotificationService
}

// NewTaskService creates a new instance of TaskService.
func NewTaskService(repo *storage.TaskRepository, projectService *ProjectService, notificationService *NotificationService) *TaskService {
	return &TaskService{
		Repo:                repo,
		ProjectService:      projectService,
		NotificationService: notificationService,
	}
}

// CreateTask handles the business logic for creating a new task and returns the hydrated object.
func (s *TaskService) CreateTask(task *models.Task, userStoryID uint, creatorID uint) (*models.Task, error) {
	task.UserStoryID = userStoryID
	task.CreatedByID = creatorID

	if err := s.Repo.CreateTask(task); err != nil {
		return nil, err
	}

	return s.Repo.GetTaskByID(task.ID)
}

// GetTaskByID retrieves a single task by its ID.
func (s *TaskService) GetTaskByID(id uint) (*models.Task, error) {
	return s.Repo.GetTaskByID(id)
}

// GetTasksByUserStoryID retrieves all tasks for a specific user story.
func (s *TaskService) GetTasksByUserStoryID(userStoryID uint) ([]models.Task, error) {
	return s.Repo.GetTasksByUserStoryID(userStoryID)
}

// UpdateTask handles the business logic for updating a task.
func (s *TaskService) UpdateTask(task *models.Task) (*models.Task, error) {
	if err := s.Repo.UpdateTask(task); err != nil {
		return nil, err
	}
	return s.Repo.GetTaskByID(task.ID)
}

// DeleteTask handles the business logic for deleting a task.
func (s *TaskService) DeleteTask(id uint) error {
	return s.Repo.DeleteTask(id)
}

// AssignTask handles the business logic for assigning a task to a user.
func (s *TaskService) AssignTask(taskID, assignToUserID uint) (*models.Task, error) {
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	if task.UserStory.ProjectID == 0 {
		return nil, fmt.Errorf("could not verify project membership: task is not linked to a project")
	}
	_, err = s.ProjectService.GetUserRoleInProject(assignToUserID, task.UserStory.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("assignment failed: user is not a member of this project")
	}

	task.AssignedToID = &assignToUserID

	updatedTask, err := s.UpdateTask(task)
	if err != nil {
		return nil, err
	}

	// --- Create Notification ---
	message := fmt.Sprintf("Se te ha asignado la tarea '%s'.", updatedTask.Title)
	link := fmt.Sprintf("/tasks/%d", updatedTask.ID)
	_, err = s.NotificationService.CreateNotification(assignToUserID, message, link)
	if err != nil {
		// Log the error but don't fail the whole operation as the main action (assignment) was successful.
		log.Printf("could not create notification for task assignment: %v", err)
	}
	// --- End Notification ---

	return updatedTask, nil
}

// UpdateTaskStatus handles the business logic for changing a task's status.
func (s *TaskService) UpdateTaskStatus(taskID uint, newStatus string) (*models.Task, error) {
	// 1. Validate the new status.
	if !models.IsValidTaskStatus(newStatus) {
		return nil, fmt.Errorf("invalid task status: %s", newStatus)
	}

	// 2. Get the task.
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	// 3. Update the status.
	task.Status = newStatus

	// 4. Save and return the updated, hydrated task.
	return s.UpdateTask(task)
}

// AddCommentToTask adds a comment to a task and notifies the assignee.
func (s *TaskService) AddCommentToTask(taskID, authorID uint, content string) (*models.TaskComment, error) {
	comment := &models.TaskComment{
		TaskID:   taskID,
		AuthorID: authorID,
		Content:  content,
	}

	if err := s.Repo.AddComment(comment); err != nil {
		return nil, err
	}

	// --- Create Notification ---
	task, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		log.Printf("could not get task for notification after commenting: %v", err)
		return comment, nil // Return the comment even if notification fails
	}

	// Only notify if there is an assignee and the assignee is not the one who commented.
	if task.AssignedToID != nil && *task.AssignedToID != authorID {
		message := fmt.Sprintf("Nuevo comentario en la tarea '%s'.", task.Title)
		link := fmt.Sprintf("/tasks/%d", task.ID)
		_, err := s.NotificationService.CreateNotification(*task.AssignedToID, message, link)
		if err != nil {
			log.Printf("could not create notification for new comment: %v", err)
		}
	}
	// --- End Notification ---

	return comment, nil
}
