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
	ActivityService     *ActivityService
}

// NewTaskService creates a new instance of TaskService.
func NewTaskService(repo *storage.TaskRepository, projectService *ProjectService, notificationService *NotificationService, activityService *ActivityService) *TaskService {
	return &TaskService{
		Repo:                repo,
		ProjectService:      projectService,
		NotificationService: notificationService,
		ActivityService:     activityService,
	}
}

// CreateTask handles the business logic for creating a new task and returns the hydrated object.
func (s *TaskService) CreateTask(task *models.Task, userStoryID uint, creatorID uint) (*models.Task, error) {
	// Keep a reference to the assigned ID, but don't save it directly on creation.
	assignedToID := task.AssignedToID
	task.AssignedToID = nil // Ensure it's nil before creation.

	task.UserStoryID = userStoryID
	task.CreatedByID = creatorID
	task.Status = models.StatusTodo // Default status

	if err := s.Repo.CreateTask(task); err != nil {
		return nil, err
	}

	// If an assignee was specified, assign the task now.
	if assignedToID != nil && *assignedToID != 0 {
		return s.AssignTask(task.ID, *assignedToID, creatorID)
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
func (s *TaskService) AssignTask(taskID, assignToUserID, assignerID uint) (*models.Task, error) {
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

	// Set the preloaded association to nil before changing the ID.
	// This forces GORM to recognize the ID change and prevents the bug.
	task.AssignedTo = nil
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

	// --- Create Activity ---
	if updatedTask.AssignedTo != nil {
		projectID, err := s.Repo.GetProjectIDForTask(taskID)
		if err == nil { // Only create activity if project ID is found
			description := fmt.Sprintf("ha asignado la tarea '%s' a %s.", updatedTask.Title, updatedTask.AssignedTo.Nombre)
			s.ActivityService.CreateActivity(
				"task_assigned",
				"task",
				taskID,
				assignerID,
				projectID,
				description,
			)
		}
	}
	// --- End Activity ---

	return updatedTask, nil
}

// UpdateTaskStatus handles the business logic for changing a task's status.
func (s *TaskService) UpdateTaskStatus(taskID uint, newStatus string, updaterID uint) (*models.Task, error) {
	// 1. Validate the new status.
	if !models.IsValidTaskStatus(newStatus) {
		return nil, fmt.Errorf("invalid task status: %s", newStatus)
	}

	// 2. Get the task to find the old status.
	originalTask, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}
	oldStatus := originalTask.Status
	newStatusTyped := models.TaskStatus(newStatus)

	// 3. If the status is the same, do nothing.
	if oldStatus == newStatusTyped {
		return originalTask, nil
	}

	// 4. Save the change and create a history record.
	// Pass the updaterID to the repository to be stored in the history.
	if err := s.Repo.UpdateTaskWithHistory(originalTask, oldStatus, newStatusTyped, updaterID); err != nil {
		return nil, err
	}

	// 5. Return the updated, hydrated task.
	updatedTask, err := s.Repo.GetTaskByID(taskID)
	if err != nil {
		return nil, err // Task should exist, but handle error just in case.
	}

	// --- Create Activity ---
	if newStatusTyped == models.StatusDone && updatedTask.AssignedToID != nil {
		projectID, err := s.Repo.GetProjectIDForTask(taskID)
		if err == nil { // Only create activity if project ID is found
			description := fmt.Sprintf("ha completado la tarea '%s'", updatedTask.Title)
			s.ActivityService.CreateActivity(
				"task_completed",
				"task",
				taskID,
				*updatedTask.AssignedToID, // The user who completed it is the assignee
				projectID,
				description,
			)
		}
	}
	// --- End Activity ---

	return updatedTask, nil
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

	// --- Create Activity ---
	projectID, err := s.Repo.GetProjectIDForTask(taskID)
	if err == nil {
		description := fmt.Sprintf("ha dejado un comentario en '%s'", task.Title)
		s.ActivityService.CreateActivity(
			"comment",
			"task",
			taskID,
			authorID,
			projectID,
			description,
		)
	}
	// --- End Activity ---

	return comment, nil
}

// GetCommentsByTaskID retrieves all comments for a specific task.
func (s *TaskService) GetCommentsByTaskID(taskID uint) ([]models.TaskComment, error) {
	// Future enhancement: Add permission check here to ensure the user can view the task.
	return s.Repo.GetCommentsByTaskID(taskID)
}

// GetTasksByUserID handles the business logic for fetching a user's tasks with filters.
func (s *TaskService) GetTasksByUserID(userID uint, projectID *uint, status *string, limit, offset int) ([]models.Task, int64, error) {
	// Business logic can be added here, e.g., validating parameters.
	// For now, it directly calls the repository.
	return s.Repo.GetTasksByUserID(userID, projectID, status, limit, offset)
}
