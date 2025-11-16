package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/utils"
	"github.com/buga/API_wrkf/websocket"
	"github.com/labstack/echo/v4"
)

// AssignTaskRequest defines the structure for assigning a task to a user.
type AssignTaskRequest struct {
	UserID uint `json:"userId" example:"2"`
}

// UpdateTaskStatusRequest defines the structure for updating a task's status.
type UpdateTaskStatusRequest struct {
	Status string `json:"status" example:"in_progress"`
}

// CreateTaskRequest defines the structure for creating a new task.
// Using a DTO with explicit json tags is more robust than binding directly to the model.
type CreateTaskRequest struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	AssignedToID   *uint     `json:"assignedToID"`
	EstimatedHours *float32  `json:"estimatedHours"`
	IsDeliverable  bool      `json:"isDeliverable"`
}

// TaskHandler handles HTTP requests for tasks.
type TaskHandler struct {
	Service     *services.TaskService
	wsManager   *websocket.WebSocketManager
	userService *services.UserService
}

// NewTaskHandler creates a new instance of TaskHandler.
func NewTaskHandler(service *services.TaskService, wsManager *websocket.WebSocketManager, userService *services.UserService) *TaskHandler {
	return &TaskHandler{
		Service:     service,
		wsManager:   wsManager,
		userService: userService,
	}
}

// CreateTask godoc
// @Summary      Create a new Task
// @Description  Creates a new task within a specific user story.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        storyId  path      int           true  "User Story ID"
// @Param        task     body      models.Task   true  "Task Details"
// @Success      201      {object}  models.Task
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/userstories/{storyId}/tasks [post]
func (h *TaskHandler) CreateTask(c echo.Context) error {
	userStoryID, err := strconv.ParseUint(c.Param("storyId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user story ID"})
	}

	req := new(CreateTaskRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	creatorID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or missing user ID from token"})
	}

	// Manually construct the task model from the DTO.
	// This is the definitive fix for the assignment-on-creation bug.
	task := &models.Task{
		Title:          req.Title,
		Description:    req.Description,
		AssignedToID:   req.AssignedToID,
		EstimatedHours: req.EstimatedHours,
		IsDeliverable:  req.IsDeliverable,
	}

	createdTask, err := h.Service.CreateTask(task, uint(userStoryID), creatorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not create task: %v", err)})
	}

	// Broadcast WebSocket event
	if h.wsManager != nil && createdTask.UserStory.ProjectID != 0 {
		h.wsManager.BroadcastTaskCreated(createdTask.UserStory.ProjectID, createdTask)
	}

	return c.JSON(http.StatusCreated, createdTask)
}

// GetTasksByUserStoryID godoc
// @Summary      Get all Tasks for a User Story
// @Description  Retrieves a list of all tasks for a specific user story.
// @Tags         Tasks
// @Produce      json
// @Param        storyId   path      int  true  "User Story ID"
// @Success      200       {array}   models.Task
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/userstories/{storyId}/tasks [get]
func (h *TaskHandler) GetTasksByUserStoryID(c echo.Context) error {
	userStoryID, err := strconv.ParseUint(c.Param("storyId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user story ID"})
	}

	tasks, err := h.Service.GetTasksByUserStoryID(uint(userStoryID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve tasks"})
	}

	return c.JSON(http.StatusOK, tasks)
}

// UpdateTask godoc
// @Summary      Update a Task
// @Description  Updates an existing task.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        taskId   path      int           true  "Task ID"
// @Param        task     body      models.Task   true  "Updated task data"
// @Success      200      {object}  models.Task
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks/{taskId} [put]
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	taskId, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	taskToUpdate, err := h.Service.GetTaskByID(uint(taskId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	if err := c.Bind(taskToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	updatedTask, err := h.Service.UpdateTask(taskToUpdate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update task"})
	}

	return c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask godoc
// @Summary      Delete a Task
// @Description  Deletes an existing task.
// @Tags         Tasks
// @Param        taskId   path      int  true  "Task ID"
// @Success      204      {object}  nil
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks/{taskId} [delete]
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	taskId, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	// Get info before deleting
	task, err := h.Service.GetTaskByID(uint(taskId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}
	projectID := task.UserStory.ProjectID

	deleterID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	deleter, err := h.userService.GetUserByID(deleterID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Delete the task
	if err := h.Service.DeleteTask(uint(taskId)); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found or could not be deleted"})
	}

	// Broadcast WebSocket event
	if h.wsManager != nil && projectID != 0 {
		h.wsManager.BroadcastTaskDeleted(projectID, uint(taskId), deleter)
	}

	return c.NoContent(http.StatusNoContent)
}

// AssignTask godoc
// @Summary      Assign a Task to a User
// @Description  Assigns an existing task to a user who is a member of the project.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        taskId   path      int                 true  "Task ID"
// @Param        assignment body      AssignTaskRequest   true  "User ID to assign the task to"
// @Success      200      {object}  models.Task
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks/{taskId}/assign [put]
func (h *TaskHandler) AssignTask(c echo.Context) error {
	taskId, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	req := new(AssignTaskRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	assignerID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	assignedTask, err := h.Service.AssignTask(uint(taskId), req.UserID, assignerID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Broadcast WebSocket event
	if h.wsManager != nil && assignedTask.UserStory.ProjectID != 0 {
		assigner, err := h.userService.GetUserByID(assignerID)
		if err == nil && assignedTask.AssignedTo != nil {
			h.wsManager.BroadcastTaskAssigned(assignedTask.UserStory.ProjectID, assignedTask.ID, assignedTask.AssignedTo, assigner)
		}
	}

	return c.JSON(http.StatusOK, assignedTask)
}

// UpdateTaskStatus godoc
// @Summary      Update a Task's Status
// @Description  Updates the status of an existing task (e.g., 'todo', 'in_progress', 'done').
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        taskId   path      int                      true  "Task ID"
// @Param        status   body      UpdateTaskStatusRequest  true  "New status for the task"
// @Success      200      {object}  models.Task
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/tasks/{taskId}/status [put]
func (h *TaskHandler) UpdateTaskStatus(c echo.Context) error {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	req := new(UpdateTaskStatusRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Get the user performing the update
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	updater, err := h.userService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Get the original task to find the old status and project ID
	originalTask, err := h.Service.GetTaskByID(uint(taskID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}
	oldStatus := originalTask.Status
	projectID := originalTask.UserStory.ProjectID

	// Update the task status, passing the updater's ID to the service layer
	updatedTask, err := h.Service.UpdateTaskStatus(uint(taskID), req.Status, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Broadcast the event
	if h.wsManager != nil {
		h.wsManager.BroadcastTaskStatusUpdated(projectID, updatedTask.ID, string(oldStatus), string(updatedTask.Status), updater)
	}

	return c.JSON(http.StatusOK, updatedTask)
}

// AddComment handles the request to add a new comment to a task.
func (h *TaskHandler) AddComment(c echo.Context) error {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	authorID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	var body struct {
		Content string `json:"content"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if body.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Comment content cannot be empty"})
	}

	comment, err := h.Service.AddCommentToTask(uint(taskID), authorID, body.Content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, comment)
}

// GetCommentsByTaskID retrieves all comments for a specific task.
func (h *TaskHandler) GetCommentsByTaskID(c echo.Context) error {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	comments, err := h.Service.GetCommentsByTaskID(uint(taskID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve comments"})
	}

	return c.JSON(http.StatusOK, comments)
}

// GetUserTasks retrieves tasks for a specific user with optional filters.
func (h *TaskHandler) GetUserTasks(c echo.Context) error {
	// 1. Get user ID from URL
	targetUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// 2. Authorization check
	requestingUserID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	userRole, _ := utils.GetUserRoleFromContext(c)

	if userRole != "admin" && requestingUserID != uint(targetUserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: You can only view your own tasks"})
	}

	// 3. Parse query parameters
	projectIDStr := c.QueryParam("project_id")
	var projectID *uint
	if projectIDStr != "" {
		pid, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil {
			p := uint(pid)
			projectID = &p
		}
	}

	statusStr := c.QueryParam("status")
	var status *string
	if statusStr != "" {
		status = &statusStr
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 50 // Default limit
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	// 4. Call the service
	tasks, total, err := h.Service.GetTasksByUserID(uint(targetUserID), projectID, status, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve tasks"})
	}

	// 5. Format the response
	type TaskResponse struct {
		ID             uint       `json:"id"`
		Title          string     `json:"title"`
		Description    string     `json:"description"`
		Status         string     `json:"status"`
		Priority       string     `json:"priority"`
		DueDate        *time.Time `json:"dueDate"`
		UserStoryID    uint       `json:"userStoryId"`
		UserStoryTitle string     `json:"userStoryTitle"`
		ProjectID      uint       `json:"projectId"`
		ProjectName    string     `json:"projectName"`
		CreatedAt      string     `json:"createdAt"`
		UpdatedAt      string     `json:"updatedAt"`
	}

	taskResponses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = TaskResponse{
			ID:             task.ID,
			Title:          task.Title,
			Description:    task.Description,
			Status:         string(task.Status),
			Priority:       task.Priority,
			DueDate:        task.DueDate,
			UserStoryID:    task.UserStoryID,
			UserStoryTitle: task.UserStory.Title,
			ProjectID:      task.UserStory.ProjectID,
			ProjectName:    task.UserStory.Project.Name,
			CreatedAt:      task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      task.UpdatedAt.Format(time.RFC3339),
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tasks":   taskResponses,
		"total":   total,
		"hasMore": total > int64(offset+len(tasks)),
	})
}
