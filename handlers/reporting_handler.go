package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// ReportingHandler handles HTTP requests for reports.
type ReportingHandler struct {
	service services.ReportingService
}

// NewReportingHandler creates a new instance of ReportingHandler.
func NewReportingHandler(service services.ReportingService) *ReportingHandler {
	return &ReportingHandler{service: service}
}

// GetProjectVelocity handles the request to get a project's velocity report.
func (h *ReportingHandler) GetProjectVelocity(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID format"})
	}

	report, err := h.service.CalculateProjectVelocity(uint(id))
	if err != nil {
		// In a real app, you might check for specific errors, e.g., gorm.ErrRecordNotFound
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}

// GetSprintCommitmentReport handles the request to get a sprint's commitment vs. completed report.
func (h *ReportingHandler) GetSprintCommitmentReport(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID format"})
	}

	report, err := h.service.CalculateSprintCommitment(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Sprint not found or could not generate report"})
	}

	return c.JSON(http.StatusOK, report)
}

// GetSprintBurndown handles the request to get a sprint's burndown chart data.
func (h *ReportingHandler) GetSprintBurndown(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID format"})
	}

	report, err := h.service.CalculateSprintBurndown(uint(id))
	if err != nil {
		// Check for the specific error we defined in the service
		if err.Error() == "sprint must have a start and end date" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		// Handle other potential errors, like sprint not found
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Sprint not found or could not generate report"})
	}

	return c.JSON(http.StatusOK, report)
}
