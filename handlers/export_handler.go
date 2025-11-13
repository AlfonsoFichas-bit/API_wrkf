package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// ExportHandler handles HTTP requests for data exportation.
type ExportHandler struct {
	Service *services.ExportService
}

// NewExportHandler creates a new instance of ExportHandler.
func NewExportHandler(service *services.ExportService) *ExportHandler {
	return &ExportHandler{Service: service}
}

// ExportProject handles the request to export a project's data to CSV.
func (h *ExportHandler) ExportProject(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	csvData, err := h.Service.ExportProjectToCSV(uint(projectID))
	if err != nil {
		// Differentiate between not found and other errors if needed
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Set headers to prompt file download
	fileName := fmt.Sprintf("project_%d_export_%s.csv", projectID, time.Now().Format("20060102"))
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	return c.Blob(http.StatusOK, "text/csv", csvData)
}
