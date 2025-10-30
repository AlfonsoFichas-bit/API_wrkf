package handlers

import (
	"net/http"
	"strconv"
	"strings"

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
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	evaluatorID, ok := c.Get("userID").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID from token"})
	}

	req := new(services.CreateEvaluationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	evaluation, err := h.Service.CreateEvaluation(uint(taskID), uint(evaluatorID), *req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "permission") || strings.Contains(err.Error(), "not a deliverable") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, evaluation)
}

// GetEvaluation handles retrieving the evaluation for a task.
func (h *EvaluationHandler) GetEvaluation(c echo.Context) error {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	evaluation, err := h.Service.GetEvaluation(uint(taskID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Evaluation not found"})
	}
	return c.JSON(http.StatusOK, evaluation)
}

// UpdateEvaluation handles updating an existing evaluation.
func (h *EvaluationHandler) UpdateEvaluation(c echo.Context) error {
	evalID, err := strconv.ParseUint(c.Param("evalId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid evaluation ID"})
	}

	evaluatorID, ok := c.Get("userID").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID from token"})
	}

	req := new(services.CreateEvaluationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	evaluation, err := h.Service.UpdateEvaluation(uint(evalID), uint(evaluatorID), *req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "permission") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, evaluation)
}
