package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/middleware"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// EvaluationHandler handles HTTP requests for evaluations.
type EvaluationHandler struct {
	Service *services.EvaluationService
}

// NewEvaluationHandler creates a new instance of EvaluationHandler.
func NewEvaluationHandler(service *services.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{Service: service}
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
	evaluatorID, err := middleware.GetUserIDFromContext(c)
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
