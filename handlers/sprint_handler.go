package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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

// CreateSprint godoc
// @Summary      Create a new Sprint
// @Description  Creates a new sprint within a specific project.
// @Tags         Sprints
// @Accept       json
// @Produce      json
// @Param        id     path      int           true  "Project ID"
// @Param        sprint body      models.Sprint true  "Sprint Details (name, goal, startDate, endDate)"
// @Success      201    {object}  models.Sprint
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/projects/{id}/sprints [post]
func (h *SprintHandler) CreateSprint(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	sprint := new(models.Sprint)
	if err := c.Bind(sprint); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	creatorID, _ := c.Get("userID").(float64)

	if err := h.Service.CreateSprint(sprint, uint(projectID), uint(creatorID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not create sprint: %v", err)})
	}

	return c.JSON(http.StatusCreated, sprint)
}

// GetSprintsByProjectID godoc
// @Summary      Get all Sprints for a project
// @Description  Retrieves a list of all sprints for a specific project.
// @Tags         Sprints
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {array}   models.Sprint
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/projects/{id}/sprints [get]
func (h *SprintHandler) GetSprintsByProjectID(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	sprints, err := h.Service.GetSprintsByProjectID(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve sprints"})
	}

	return c.JSON(http.StatusOK, sprints)
}

// GetSprintByID godoc
// @Summary      Get a single Sprint
// @Description  Retrieves details of a single sprint by its ID.
// @Tags         Sprints
// @Produce      json
// @Param        sprintId   path      int  true  "Sprint ID"
// @Success      200  {object}  models.Sprint
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/sprints/{sprintId} [get]
func (h *SprintHandler) GetSprintByID(c echo.Context) error {
	sprintID, err := strconv.ParseUint(c.Param("sprintId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID"})
	}

	sprint, err := h.Service.GetSprintByID(uint(sprintID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Sprint not found"})
	}

	return c.JSON(http.StatusOK, sprint)
}

// UpdateSprint godoc
// @Summary      Update a Sprint
// @Description  Updates an existing sprint.
// @Tags         Sprints
// @Accept       json
// @Produce      json
// @Param        sprintId   path      int           true  "Sprint ID"
// @Param        sprint     body      models.Sprint true  "Updated sprint data"
// @Success      200        {object}  models.Sprint
// @Failure      400        {object}  map[string]string
// @Failure      401        {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/sprints/{sprintId} [put]
func (h *SprintHandler) UpdateSprint(c echo.Context) error {
	sprintID, err := strconv.ParseUint(c.Param("sprintId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID"})
	}

	sprintToUpdate, err := h.Service.GetSprintByID(uint(sprintID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Sprint not found"})
	}

	if err := c.Bind(sprintToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := h.Service.UpdateSprint(sprintToUpdate); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update sprint"})
	}

	return c.JSON(http.StatusOK, sprintToUpdate)
}

// DeleteSprint godoc
// @Summary      Delete a Sprint
// @Description  Deletes an existing sprint.
// @Tags         Sprints
// @Param        sprintId   path      int  true  "Sprint ID"
// @Success      204        {object}  nil
// @Failure      401        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/sprints/{sprintId} [delete]
func (h *SprintHandler) DeleteSprint(c echo.Context) error {
	sprintID, err := strconv.ParseUint(c.Param("sprintId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID"})
	}

	if err := h.Service.DeleteSprint(uint(sprintID)); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Sprint not found or could not be deleted"})
	}

	return c.NoContent(http.StatusNoContent)
}
