package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// BurndownHandler handles HTTP requests related to burndown charts.
type BurndownHandler struct {
	service *services.BurndownService
}

// NewBurndownHandler creates a new instance of BurndownHandler.
func NewBurndownHandler(service *services.BurndownService) *BurndownHandler {
	return &BurndownHandler{service: service}
}

// GetBurndownChart handles the request to generate and return a burndown chart.
func (h *BurndownHandler) GetBurndownChart(c echo.Context) error {
	sprintIDStr := c.Param("id")
	sprintID, err := strconv.ParseUint(sprintIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID"})
	}

	chart, err := h.service.GenerateBurndownChart(uint(sprintID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, chart)
}
