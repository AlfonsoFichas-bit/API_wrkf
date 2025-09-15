
package handlers

import (
	"API_wrkf/models"
	"API_wrkf/services"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// AddMemberRequest defines the structure for a request to add a member to a project.
type AddMemberRequest struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
}

// ProjectHandler handles HTTP requests for projects.
type ProjectHandler struct {
	Service *services.ProjectService
}

// NewProjectHandler creates a new instance of ProjectHandler.
func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{Service: service}
}

// CreateProject handles the HTTP request to create a new project.
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	project := new(models.Project)
	if err := c.Bind(project); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userID, _ := c.Get("userID").(float64)

	if err := h.Service.CreateProject(project, uint(userID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not create project: %v", err)})
	}

	return c.JSON(http.StatusCreated, project)
}

// AddMemberToProject handles the HTTP request to add a member to a project.
func (h *ProjectHandler) AddMemberToProject(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	req := new(AddMemberRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	member, err := h.Service.AddMemberToProject(uint(projectID), req.UserID, req.Role)
	if err != nil {
		if strings.Contains(err.Error(), "invalid project role") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not add member: %v", err)})
	}

	return c.JSON(http.StatusCreated, member)
}

// GetAllProjects handles the HTTP request to retrieve all projects.
func (h *ProjectHandler) GetAllProjects(c echo.Context) error {
	projects, err := h.Service.GetProjects()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve projects"})
	}
	return c.JSON(http.StatusOK, projects)
}

// GetProjectByID handles the HTTP request to retrieve a single project by its ID.
func (h *ProjectHandler) GetProjectByID(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	project, err := h.Service.GetProjectByID(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
	}

	return c.JSON(http.StatusOK, project)
}

// UpdateProject handles the HTTP request to update a project.
func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, _ := c.Get("userID").(float64)
	userRole, _ := c.Get("userRole").(string)

	updatedProject, err := h.Service.UpdateProject(uint(projectID), updates, uint(userID), userRole)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedProject)
}

// DeleteProject handles the HTTP request to delete a project.
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	userID, _ := c.Get("userID").(float64)
	userRole, _ := c.Get("userRole").(string)

	if err := h.Service.DeleteProject(uint(projectID), uint(userID), userRole); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
