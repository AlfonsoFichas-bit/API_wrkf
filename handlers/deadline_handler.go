package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/utils"
	"github.com/labstack/echo/v4"
)

// DeadlineHandler handles HTTP requests for deadlines.
type DeadlineHandler struct {
	Service     *services.DeadlineService
	TeamService *services.TeamService
}

// NewDeadlineHandler creates a new instance of DeadlineHandler.
func NewDeadlineHandler(service *services.DeadlineService, teamService *services.TeamService) *DeadlineHandler {
	return &DeadlineHandler{Service: service, TeamService: teamService}
}

// GetUpcomingDeadlines handles the request to get upcoming deadlines.
func (h *DeadlineHandler) GetUpcomingDeadlines(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	userRole, _ := utils.GetUserRoleFromContext(c)

	// Parse query params
	days, _ := strconv.Atoi(c.QueryParam("days"))
	if days <= 0 || days > 90 {
		days = 30 // Default days
	}

	projectIDStr := c.QueryParam("project_id")
	var projectID *uint
	if projectIDStr != "" {
		pid, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil {
			p := uint(pid)
			projectID = &p
		}
	}

	deadlineTypeStr := c.QueryParam("type")
	var deadlineType *string
	if deadlineTypeStr != "" {
		deadlineType = &deadlineTypeStr
	}

	deadlines, err := h.Service.GetUpcomingDeadlines(userID, userRole, days, projectID, deadlineType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Format response
	type DeadlineResponse struct {
		ID            uint      `json:"id"`
		Title         string    `json:"title"`
		Type          string    `json:"type"`
		ProjectName   *string   `json:"projectName"`
		ProjectID     *uint     `json:"projectId"`
		Date          time.Time `json:"date"`
		IsUrgent      bool      `json:"isUrgent"`
		DaysRemaining int       `json:"daysRemaining"`
		AffectedTeams []string  `json:"affectedTeams"`
		Description   *string   `json:"description,omitempty"` // Assuming description can be null
	}

	response := make([]DeadlineResponse, len(deadlines))
	urgentCount := 0
	now := time.Now()

	for i, d := range deadlines {
		daysRemaining := int(math.Ceil(d.Date.Sub(now).Hours() / 24))
		isUrgent := daysRemaining <= 3
		if isUrgent {
			urgentCount++
		}

		var projectName *string
		affectedTeams := []string{}
		if d.ProjectID != nil {
			if d.Project != nil {
				projectName = &d.Project.Name
			}
			// If it's a sprint end, assume it affects all teams in the project
			if d.Type == "sprint_end" {
				teams, err := h.TeamService.GetTeamsByProjectID(*d.ProjectID)
				if err == nil {
					for _, team := range teams {
						affectedTeams = append(affectedTeams, team.Name)
					}
				}
			}
		}

		response[i] = DeadlineResponse{
			ID:            d.ID,
			Title:         d.Title,
			Type:          d.Type,
			ProjectName:   projectName,
			ProjectID:     d.ProjectID,
			Date:          d.Date,
			IsUrgent:      isUrgent,
			DaysRemaining: daysRemaining,
			AffectedTeams: affectedTeams,
			// Description would need to be added to the model if required
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"deadlines":   response,
		"total":       len(response),
		"urgentCount": urgentCount,
	})
}
