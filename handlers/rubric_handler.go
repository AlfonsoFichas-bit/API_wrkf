package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

type RubricHandler struct {
	service services.RubricService
}

func NewRubricHandler(service services.RubricService) *RubricHandler {
	return &RubricHandler{service: service}
}

// GetAllRubrics handles GET requests to fetch all rubrics.
// It supports filtering via query parameters like "isTemplate".
func (h *RubricHandler) GetAllRubrics(c echo.Context) error {
	filters := make(map[string]interface{})
	if isTemplate := c.QueryParam("isTemplate"); isTemplate != "" {
		filters["is_template"] = isTemplate
	}
	// Note: userId and projectId filters from docs would be added here
	// once auth and project association are fully integrated.

	rubrics, err := h.service.GetAllRubrics(filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rubrics)
}

// GetRubricByID handles GET requests for a single rubric by its ID.
func (h *RubricHandler) GetRubricByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	rubric, err := h.service.GetRubricByID(uint(id))
	if err != nil {
		// A common pattern is to check for gorm.ErrRecordNotFound specifically
		// but for now, a general error is fine.
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Rubric not found"})
	}
	return c.JSON(http.StatusOK, rubric)
}

// CreateRubric handles POST requests to create a new rubric.
func (h *RubricHandler) CreateRubric(c echo.Context) error {
	var rubric models.Rubric
	if err := c.Bind(&rubric); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.service.CreateRubric(&rubric); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, rubric)
}

// UpdateRubric handles PUT requests to update an existing rubric.
func (h *RubricHandler) UpdateRubric(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	var rubric models.Rubric
	if err := c.Bind(&rubric); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	rubric.ID = uint(id) // Ensure the ID from the URL is used

	if err := h.service.UpdateRubric(&rubric); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rubric)
}

// DeleteRubric handles DELETE requests to remove a rubric.
func (h *RubricHandler) DeleteRubric(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	if err := h.service.DeleteRubric(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// DuplicateRubric handles POST requests to duplicate a rubric.
func (h *RubricHandler) DuplicateRubric(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	newRubric, err := h.service.DuplicateRubric(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, newRubric)
}
