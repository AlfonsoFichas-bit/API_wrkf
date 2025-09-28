package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// MetricHandler handles HTTP requests for metrics.
type MetricHandler struct {
	metricService *services.MetricService
}

// NewMetricHandler creates a new instance of MetricHandler.
func NewMetricHandler(metricService *services.MetricService) *MetricHandler {
	return &MetricHandler{metricService: metricService}
}

// GetBurndownChart retrieves the burndown chart for a given sprint.
// @Summary Get a sprint's burndown chart
// @Description Get a sprint's burndown chart
// @Tags Metrics
// @Accept json
// @Produce json
// @Param sprintId path int true "Sprint ID"
// @Success 200 {array} services.BurndownChartPoint
// @Router /api/sprints/{sprintId}/metrics/burndown [get]
func (h *MetricHandler) GetBurndownChart(c echo.Context) error {
	sprintID, err := strconv.Atoi(c.Param("sprintId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid sprint ID")
	}

	burndownChart, err := h.metricService.GetBurndownChart(uint(sprintID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, burndownChart)
}

// GetTeamVelocity retrieves the team velocity for a given project.
// @Summary Get a project's team velocity
// @Description Get a project's team velocity
// @Tags Metrics
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Success 200 {array} services.TeamVelocity
// @Router /api/projects/{projectId}/metrics/velocity [get]
func (h *MetricHandler) GetTeamVelocity(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid project ID")
	}

	teamVelocity, err := h.metricService.GetTeamVelocity(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, teamVelocity)
}

// GetWorkDistribution retrieves the work distribution for a given sprint.
// @Summary Get a sprint's work distribution
// @Description Get a sprint's work distribution
// @Tags Metrics
// @Accept json
// @Produce json
// @Param sprintId path int true "Sprint ID"
// @Success 200 {array} services.WorkDistribution
// @Router /api/sprints/{sprintId}/metrics/work-distribution [get]
func (h *MetricHandler) GetWorkDistribution(c echo.Context) error {
	sprintID, err := strconv.Atoi(c.Param("sprintId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid sprint ID")
	}

	workDistribution, err := h.metricService.GetWorkDistribution(uint(sprintID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, workDistribution)
}