package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

type EvaluationHandler struct {
	service services.IEvaluationService
}

func NewEvaluationHandler(service services.IEvaluationService) *EvaluationHandler {
	return &EvaluationHandler{service: service}
}

func (h *EvaluationHandler) CreateEvaluation(c echo.Context) error {
	var evaluation models.Evaluation
	if err := c.Bind(&evaluation); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.service.CreateEvaluation(&evaluation); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, evaluation)
}

func (h *EvaluationHandler) GetEvaluation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 0 {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	evaluation, err := h.service.GetEvaluationByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, "Evaluation not found")
	}

	return c.JSON(http.StatusOK, evaluation)
}

func (h *EvaluationHandler) GetStudentEvaluations(c echo.Context) error {
	studentID, err := strconv.Atoi(c.Param("studentId"))
	if err != nil || studentID < 0 {
		return c.JSON(http.StatusBadRequest, "Invalid Student ID")
	}

	evaluations, err := h.service.GetEvaluationsByStudent(uint(studentID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, evaluations)
}
