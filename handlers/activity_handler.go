package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/utils"
	"github.com/labstack/echo/v4"
)

// ActivityHandler handles HTTP requests for activities.
type ActivityHandler struct {
	Service *services.ActivityService
}

// NewActivityHandler creates a new instance of ActivityHandler.
func NewActivityHandler(service *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{Service: service}
}

// GetRecentActivities handles the request to get recent activities.
func (h *ActivityHandler) GetRecentActivities(c echo.Context) error {
	requestingUserID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Parse query params
	projectIDStr := c.QueryParam("project_id")
	var projectID *uint
	if projectIDStr != "" {
		pid, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil {
			p := uint(pid)
			projectID = &p
		}
	}

	filterUserIDStr := c.QueryParam("user_id")
	var filterUserID *uint
	if filterUserIDStr != "" {
		uid, err := strconv.ParseUint(filterUserIDStr, 10, 32)
		if err == nil {
			u := uint(uid)
			filterUserID = &u
		}
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10 // Default limit
	}

	activities, err := h.Service.GetRecentActivities(requestingUserID, projectID, filterUserID, limit)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	// Format response
	type ActivityResponse struct {
		ID          uint      `json:"id"`
		Type        string    `json:"type"`
		UserName    string    `json:"userName"`
		UserID      uint      `json:"userId"`
		Description string    `json:"description"`
		EntityType  string    `json:"entityType"`
		EntityID    uint      `json:"entityId"`
		ProjectID   uint      `json:"projectId"`
		ProjectName string    `json:"projectName"`
		Timestamp   time.Time `json:"timestamp"`
	}

	response := make([]ActivityResponse, len(activities))
	for i, a := range activities {
		response[i] = ActivityResponse{
			ID:          a.ID,
			Type:        a.Type,
			UserName:    a.User.Nombre + " " + a.User.ApellidoPaterno,
			UserID:      a.UserID,
			Description: a.Description,
			EntityType:  a.EntityType,
			EntityID:    a.EntityID,
			ProjectID:   a.ProjectID,
			ProjectName: a.Project.Name,
			Timestamp:   a.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"activities": response,
		"total":      len(response),
	})
}
