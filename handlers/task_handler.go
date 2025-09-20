package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"

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

// TaskHandler handles HTTP requests for tasks.
type TaskHandler struct {
	Service *services.TaskService
}

// NewTaskHandler creates a new instance of TaskHandler.
func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
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

	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	creatorID, ok := c.Get("userID").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or missing user ID from token"})
	}

	createdTask, err := h.Service.CreateTask(task, uint(userStoryID), uint(creatorID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not create task: %v", err)})
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

	if err := h.Service.DeleteTask(uint(taskId)); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found or could not be deleted"})
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

	assignedTask, err := h.Service.AssignTask(uint(taskId), req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
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
	taskId, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	req := new(UpdateTaskStatusRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updatedTask, err := h.Service.UpdateTaskStatus(uint(taskId), req.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedTask)
}
