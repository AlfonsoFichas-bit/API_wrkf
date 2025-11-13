package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/buga/API_wrkf/middleware"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/services"
	"github.com/labstack/echo/v4"
)

// EventHandler handles HTTP requests for calendar events.
type EventHandler struct {
	Service *services.EventService
}

// NewEventHandler creates a new instance of EventHandler.
func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{Service: service}
}

// CreateEvent handles the creation of a new event in a project.
func (h *EventHandler) CreateEvent(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	var event models.Event
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	creatorID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Could not get user from token"})
	}

	createdEvent, err := h.Service.CreateEvent(&event, uint(projectID), creatorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdEvent)
}

// GetEvents handles fetching events for a project, optionally filtering by a date range.
func (h *EventHandler) GetEvents(c echo.Context) error {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Could not get user from token"})
	}

	// Parse date range from query params
	startStr := c.QueryParam("start") // e.g., "2023-01-01"
	endStr := c.QueryParam("end")     // e.g., "2023-01-31"

	// Default to a wide range if not provided
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		start = time.Now().AddDate(-1, 0, 0) // Default to one year ago
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		end = time.Now().AddDate(1, 0, 0) // Default to one year from now
	}

	events, err := h.Service.GetEventsForProject(uint(projectID), userID, start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, events)
}

// UpdateEvent handles updating an existing event.
func (h *EventHandler) UpdateEvent(c echo.Context) error {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event ID"})
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Could not get user from token"})
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	updatedEvent, err := h.Service.UpdateEvent(uint(eventID), userID, updates)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent handles deleting an event.
func (h *EventHandler) DeleteEvent(c echo.Context) error {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event ID"})
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Could not get user from token"})
	}

	if err := h.Service.DeleteEvent(uint(eventID), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
