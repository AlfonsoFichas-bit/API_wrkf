package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"

	"github.com/labstack/echo/v4"
)

// AssignStoryRequest defines the structure for assigning a user story to a sprint.
type AssignStoryRequest struct {
	UserStoryID uint `json:"userStoryId" example:"1"`
}

// UserStoryHandler handles HTTP requests for user stories.
type UserStoryHandler struct {
	Service *services.UserStoryService
}

// NewUserStoryHandler creates a new instance of UserStoryHandler.
func NewUserStoryHandler(service *services.UserStoryService) *UserStoryHandler {
	return &UserStoryHandler{Service: service}
}

// CreateUserStory godoc
// @Summary      Create a new User Story
// @Description  Creates a new user story within a specific project.
// @Tags         User Stories
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Project ID"
// @Param        userStory  body      models.UserStory  true  "User Story Details"
// @Success      201      {object}  models.UserStory
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/projects/{id}/userstories [post]
func (h *UserStoryHandler) CreateUserStory(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	userStory := new(models.UserStory)
	if err := c.Bind(userStory); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	creatorID, ok := c.Get("userID").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or missing user ID from token"})
	}

	if err := h.Service.CreateUserStory(userStory, uint(projectID), uint(creatorID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not create user story: %v", err)})
	}

	return c.JSON(http.StatusCreated, userStory)
}

// GetUserStoriesByProjectID godoc
// @Summary      Get all User Stories for a project
// @Description  Retrieves a list of all user stories (the Product Backlog) for a specific project.
// @Tags         User Stories
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {array}   models.UserStory
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/projects/{id}/userstories [get]
func (h *UserStoryHandler) GetUserStoriesByProjectID(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	userStories, err := h.Service.GetUserStoriesByProjectID(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve user stories"})
	}

	return c.JSON(http.StatusOK, userStories)
}

// GetUserStoryByID godoc
// @Summary      Get a single User Story
// @Description  Retrieves details of a single user story by its ID, including related project, sprint, and user data.
// @Tags         User Stories
// @Produce      json
// @Param        storyId   path      int  true  "User Story ID"
// @Success      200       {object}  models.UserStory
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/userstories/{storyId} [get]
func (h *UserStoryHandler) GetUserStoryByID(c echo.Context) error {
	storyID, err := strconv.ParseUint(c.Param("storyId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user story ID"})
	}

	userStory, err := h.Service.GetUserStoryByID(uint(storyID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User story not found"})
	}

	return c.JSON(http.StatusOK, userStory)
}

// UpdateUserStory godoc
// @Summary      Update a User Story
// @Description  Updates an existing user story. Requires admin, product owner, or scrum master role.
// @Tags         User Stories
// @Accept       json
// @Produce      json
// @Param        storyId  path      int  true  "User Story ID"
// @Param        updates  body      map[string]interface{}  true  "Fields to update"
// @Success      200      {object}  models.UserStory
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/userstories/{storyId} [put]
func (h *UserStoryHandler) UpdateUserStory(c echo.Context) error {
	storyID, err := strconv.ParseUint(c.Param("storyId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user story ID"})
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, _ := c.Get("userID").(float64)
	platformRole, _ := c.Get("userRole").(string)

	updatedStory, err := h.Service.UpdateUserStory(uint(storyID), uint(userID), platformRole, updates)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedStory)
}

// UnassignUserStoryFromSprint godoc
// @Summary      Unassign a User Story from a Sprint
// @Description  Unassigns a user story from a sprint.
// @Tags         Sprints
// @Accept       json
// @Produce      json
// @Param        sprintId  path      int  true  "Sprint ID"
// @Param        storyId   path      int  true  "User Story ID"
// @Success      200       {object}  models.UserStory
// @Failure      400       {object}  map[string]string
// @Failure      403       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/sprints/{sprintId}/userstories/{storyId} [delete]
func (h *UserStoryHandler) UnassignUserStoryFromSprint(c echo.Context) error {
	sprintID, err := strconv.ParseUint(c.Param("sprintId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID"})
	}

	storyID, err := strconv.ParseUint(c.Param("storyId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user story ID"})
	}

	userID, _ := c.Get("userID").(float64)
	platformRole, _ := c.Get("userRole").(string)

	updatedStory, err := h.Service.UnassignUserStoryFromSprint(uint(sprintID), uint(storyID), uint(userID), platformRole)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedStory)
}

// DeleteUserStory godoc
// @Summary      Delete a User Story
// @Description  Deletes an existing user story. Requires admin, product owner, or scrum master role.
// @Tags         User Stories
// @Param        storyId  path      int  true  "User Story ID"
// @Success      204      {object}  nil
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/userstories/{storyId} [delete]
func (h *UserStoryHandler) DeleteUserStory(c echo.Context) error {
	storyID, err := strconv.ParseUint(c.Param("storyId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user story ID"})
	}

	userID, _ := c.Get("userID").(float64)
	platformRole, _ := c.Get("userRole").(string)

	if err := h.Service.DeleteUserStory(uint(storyID), uint(userID), platformRole); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// AssignUserStoryToSprint godoc
// @Summary      Assign a User Story to a Sprint
// @Description  Assigns an existing user story to a sprint. Requires admin, product owner, or scrum master role.
// @Tags         Sprints
// @Accept       json
// @Produce      json
// @Param        sprintId  path      int                 true  "Sprint ID"
// @Param        assignment  body      AssignStoryRequest  true  "User Story ID to assign"
// @Success      200       {object}  models.UserStory
// @Failure      400       {object}  map[string]string
// @Failure      403       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/sprints/{sprintId}/userstories [post]
func (h *UserStoryHandler) AssignUserStoryToSprint(c echo.Context) error {
	sprintID, err := strconv.ParseUint(c.Param("sprintId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sprint ID"})
	}

	req := new(AssignStoryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, _ := c.Get("userID").(float64)
	platformRole, _ := c.Get("userRole").(string)

	updatedStory, err := h.Service.AssignUserStoryToSprint(uint(sprintID), req.UserStoryID, uint(userID), platformRole)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedStory)
}
