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
