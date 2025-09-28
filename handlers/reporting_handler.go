package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// ReportingHandler handles HTTP requests for reports.
type ReportingHandler struct {
	reportingService *services.ReportingService
}

// NewReportingHandler creates a new instance of ReportingHandler.
func NewReportingHandler(reportingService *services.ReportingService) *ReportingHandler {
	return &ReportingHandler{reportingService: reportingService}
}

// GenerateProjectReport generates a report for a given project.
// @Summary Generate a project report
// @Description Generate a project report
// @Tags Reports
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Success 200 {object} models.Report
// @Router /api/projects/{projectId}/reports/generate [post]
func (h *ReportingHandler) GenerateProjectReport(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid project ID")
	}

	report, err := h.reportingService.GenerateProjectReport(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, report)
}