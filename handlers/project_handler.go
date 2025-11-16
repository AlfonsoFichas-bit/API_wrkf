package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"

	"github.com/labstack/echo/v4"
)

// AddMemberRequest define la estructura de una solicitud para añadir un miembro a un proyecto.
type AddMemberRequest struct {
	UserID uint   `json:"userId" example:"2"`
	Role   string `json:"role" example:"team_developer"`
}

// ProjectHandler gestiona las solicitudes HTTP para projects.
type ProjectHandler struct {
	Service *services.ProjectService
}

// NewProjectHandler crea una nueva instancia de ProjectHandler.
func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{Service: service}
}

// CreateProject godoc
// @Summary      Create a new project
// @Description  Creates a new project for the authenticated user.
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        project  body      models.Project  true  "Project details"
// @Success      201      {object}  models.Project
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/projects [post]
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

// AddMemberRequest define la estructura de una solicitud para añadir un miembro a un proyecto.
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
		if strings.Contains(err.Error(), "user is already a member") {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "invalid project role") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Could not add member: %v", err)})
	}

	return c.JSON(http.StatusCreated, member)
}

// GetAllProjects gestiona la solicitud HTTP para recuperar todos los proyectos.
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

// GetUnassignedUsers godoc
// @Summary      Get Unassigned Users
// @Description  Retrieves a list of users who are not admins and are not already members of a specific project.
// @Tags         Projects
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {array}   models.User
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/projects/{id}/unassigned-users [get]
func (h *ProjectHandler) GetUnassignedUsers(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	users, err := h.Service.GetUnassignedUsers(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve unassigned users"})
	}

	return c.JSON(http.StatusOK, users)
}

// GetProjectMembers handles the HTTP request to retrieve all members of a project with enhanced details.
func (h *ProjectHandler) GetProjectMembers(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	members, teams, err := h.Service.GetProjectMembers(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve project members"})
	}

	// Format the response according to the new specification
	type MemberResponse struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		Role        string    `json:"role"` // Platform role
		ProjectRole string    `json:"projectRole"`
		TeamName    *string   `json:"teamName,omitempty"`
		TeamID      *uint     `json:"teamId,omitempty"`
		JoinedAt    time.Time `json:"joinedAt"`
		IsActive    bool      `json:"isActive"`
		// Avatar field would be added here if available in the User model
	}

	memberResponses := make([]MemberResponse, len(members))
	for i, m := range members {
		var teamName *string
		if m.Team != nil {
			teamName = &m.Team.Name
		}

		memberResponses[i] = MemberResponse{
			ID:          m.User.ID,
			Name:        m.User.Nombre + " " + m.User.ApellidoPaterno,
			Email:       m.User.Correo,
			Role:        m.User.Role,
			ProjectRole: m.Role,
			TeamName:    teamName,
			TeamID:      m.TeamID,
			JoinedAt:    m.CreatedAt,
			IsActive:    true, // Assuming all returned members are active
		}
	}

	type TeamResponse struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		MemberCount int    `json:"memberCount"`
	}

	teamResponses := make([]TeamResponse, len(teams))
	for i, t := range teams {
		teamResponses[i] = TeamResponse{
			ID:          t.ID,
			Name:        t.Name,
			MemberCount: len(t.Members),
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"members": memberResponses,
		"total":   len(memberResponses),
		"teams":   teamResponses,
	})
}

// GetActiveSprint godoc
// @Summary      Get Active Sprint
// @Description  Retrieves the currently active sprint for a project.
// @Tags         Projects
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  models.Sprint
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /api/projects/{id}/active-sprint [get]
func (h *ProjectHandler) GetActiveSprint(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	sprint, err := h.Service.GetActiveSprint(uint(projectID))
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No active sprint found for this project"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve active sprint"})
	}

	return c.JSON(http.StatusOK, sprint)
}
