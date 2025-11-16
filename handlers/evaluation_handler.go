package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/storage"
	"github.com/buga/API_wrkf/utils"
	"github.com/labstack/echo/v4"
)

// EvaluationHandler handles HTTP requests for evaluations.
type EvaluationHandler struct {
	Service       *services.EvaluationService
	ProjectRepo *storage.ProjectRepository
}

// NewEvaluationHandler creates a new instance of EvaluationHandler.
func NewEvaluationHandler(service *services.EvaluationService, projectRepo *storage.ProjectRepository) *EvaluationHandler {
	return &EvaluationHandler{Service: service, ProjectRepo: projectRepo}
}

// CreateEvaluation handles the creation of a new evaluation for a task.
func (h *EvaluationHandler) CreateEvaluation(c echo.Context) error {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid task ID"})
	}

	var req services.CreateEvaluationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	// Get the evaluator's ID from the JWT token via the context.
	evaluatorID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Could not get user from token"})
	}

	evaluation, err := h.Service.CreateEvaluation(uint(taskID), evaluatorID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, evaluation)
}

// GetEvaluationsByTaskID handles fetching all evaluations for a specific task.
func (h *EvaluationHandler) GetEvaluationsByTaskID(c echo.Context) error {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid task ID"})
	}

	evaluations, err := h.Service.GetEvaluationsByTaskID(uint(taskID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if len(evaluations) == 0 {
		return c.JSON(http.StatusOK, []interface{}{}) // Return empty array instead of null
	}

	return c.JSON(http.StatusOK, evaluations)
}

// GetPendingEvaluations handles fetching tasks pending evaluation for a teacher.
func (h *EvaluationHandler) GetPendingEvaluations(c echo.Context) error {
	teacherID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Parse query params
	projectIDStr := c.QueryParam("project_id")
	var projectID *uint
	if projectIDStr != "" {
		pid, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil {
			p := uint(pid)
			projectID = &p
		}
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20 // Default limit
	}

	tasks, total, err := h.Service.GetPendingEvaluations(teacherID, projectID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Format response
	type SubmittedByResponse struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type EvaluationResponse struct {
		TaskID      uint                `json:"taskId"`
		Title       string              `json:"title"`
		Description string              `json:"description"`
		TeamName    *string             `json:"teamName,omitempty"`
		ProjectID   uint                `json:"projectId"`
		ProjectName string              `json:"projectName"`
		SubmittedAt *time.Time          `json:"submittedAt"`
		SubmittedBy SubmittedByResponse `json:"submittedBy"`
		Urgency     string              `json:"urgency"`
	}

	response := make([]EvaluationResponse, len(tasks))
	urgentCount := 0
	for i, task := range tasks {
		urgency := "low"
		if task.SubmittedAt != nil && time.Since(*task.SubmittedAt).Hours() > 72 {
			urgency = "high"
			urgentCount++
		} else if task.SubmittedAt != nil && time.Since(*task.SubmittedAt).Hours() > 24 {
			urgency = "medium"
		}

		var submittedBy SubmittedByResponse
		var teamName *string
		if task.SubmittedBy != nil {
			submittedBy = SubmittedByResponse{
				ID:    task.SubmittedBy.ID,
				Name:  task.SubmittedBy.Nombre + " " + task.SubmittedBy.ApellidoPaterno,
				Email: task.SubmittedBy.Correo,
			}
			// Get team name
			member, err := h.ProjectRepo.GetProjectMemberByUserID(task.SubmittedBy.ID, task.UserStory.ProjectID)
			if err == nil && member.Team != nil {
				teamName = &member.Team.Name
			}
		}

		response[i] = EvaluationResponse{
			TaskID:      task.ID,
			Title:       task.Title,
			Description: task.Description,
			TeamName:    teamName,
			ProjectID:   task.UserStory.ProjectID,
			ProjectName: task.UserStory.Project.Name,
			SubmittedAt: task.SubmittedAt,
			SubmittedBy: submittedBy,
			Urgency:     urgency,
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"evaluations": response,
		"total":       total,
		"urgentCount": urgentCount,
	})
}
