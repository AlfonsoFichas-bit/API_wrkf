
package handlers

import (
	"net/http"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"

	"github.com/labstack/echo/v4"
)

// SprintHandler handles HTTP requests for sprints.
type SprintHandler struct {
	Service *services.SprintService
}

// NewSprintHandler creates a new instance of SprintHandler.
func NewSprintHandler(service *services.SprintService) *SprintHandler {
	return &SprintHandler{Service: service}
}

// CreateSprint handles the HTTP request to create a new sprint.
func (h *SprintHandler) CreateSprint(c echo.Context) error {
	sprint := new(models.Sprint)
	if err := c.Bind(sprint); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Assuming the creator's ID is needed, similar to projects
	userID, _ := c.Get("userID").(float64)

	if err := h.Service.CreateSprint(sprint, uint(userID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create sprint"})
	}

	return c.JSON(http.StatusCreated, sprint)
}
