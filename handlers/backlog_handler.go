package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// BacklogHandler handles HTTP requests for the backlog.
type BacklogHandler struct {
	backlogService *services.BacklogService
}

// NewBacklogHandler creates a new instance of BacklogHandler.
func NewBacklogHandler(backlogService *services.BacklogService) *BacklogHandler {
	return &BacklogHandler{backlogService: backlogService}
}

// GetProductBacklog retrieves the product backlog for a given project.
// @Summary Get a project's product backlog
// @Description Get a project's product backlog
// @Tags Backlog
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Success 200 {array} models.UserStory
// @Router /api/projects/{projectId}/backlog [get]
func (h *BacklogHandler) GetProductBacklog(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid project ID")
	}

	backlog, err := h.backlogService.GetProductBacklog(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, backlog)
}

// UpdateUserStoryStatus updates the status of a user story.
// @Summary Update a user story's status
// @Description Update a user story's status
// @Tags Backlog
// @Accept json
// @Produce json
// @Param storyId path int true "User Story ID"
// @Param status body string true "New status"
// @Success 200 {object} models.UserStory
// @Router /api/userstories/{storyId}/status [put]
func (h *BacklogHandler) UpdateUserStoryStatus(c echo.Context) error {
	storyID, err := strconv.Atoi(c.Param("storyId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user story ID")
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.Bind(&statusUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.backlogService.UpdateUserStoryStatus(uint(storyID), statusUpdate.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "User story status updated successfully")
}