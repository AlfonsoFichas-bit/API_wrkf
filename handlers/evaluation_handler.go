package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// EvaluationHandler handles HTTP requests for evaluations.
type EvaluationHandler struct {
	evaluationService *services.EvaluationService
}

// NewEvaluationHandler creates a new instance of EvaluationHandler.
func NewEvaluationHandler(evaluationService *services.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{evaluationService: evaluationService}
}

// CreateEvaluation creates a new evaluation.
// @Summary Create a new evaluation
// @Description Create a new evaluation
// @Tags Evaluations
// @Accept json
// @Produce json
// @Param evaluation body models.Evaluation true "Evaluation object"
// @Success 201 {object} models.Evaluation
// @Router /api/evaluations [post]
func (h *EvaluationHandler) CreateEvaluation(c echo.Context) error {
	evaluation := new(models.Evaluation)
	if err := c.Bind(evaluation); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.evaluationService.CreateEvaluation(evaluation); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, evaluation)
}

// GetEvaluationByID retrieves a single evaluation by its ID.
// @Summary Get an evaluation by ID
// @Description Get an evaluation by ID
// @Tags Evaluations
// @Accept json
// @Produce json
// @Param id path int true "Evaluation ID"
// @Success 200 {object} models.Evaluation
// @Router /api/evaluations/{id} [get]
func (h *EvaluationHandler) GetEvaluationByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid evaluation ID")
	}

	evaluation, err := h.evaluationService.GetEvaluationByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, evaluation)
}

// GetEvaluationsByStudentID retrieves all evaluations for a given student ID.
// @Summary Get all evaluations for a student
// @Description Get all evaluations for a student
// @Tags Evaluations
// @Accept json
// @Produce json
// @Param studentId path int true "Student ID"
// @Success 200 {array} models.Evaluation
// @Router /api/students/{studentId}/evaluations [get]
func (h *EvaluationHandler) GetEvaluationsByStudentID(c echo.Context) error {
	studentID, err := strconv.Atoi(c.Param("studentId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid student ID")
	}

	evaluations, err := h.evaluationService.GetEvaluationsByStudentID(uint(studentID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, evaluations)
}