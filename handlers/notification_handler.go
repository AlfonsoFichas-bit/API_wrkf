package handlers

import (
	"net/http"
	"strconv"

	"github.com/buga/API_wrkf/services"
	"github.com/buga/API_wrkf/utils"
	"github.com/labstack/echo/v4"
)

// NotificationHandler handles HTTP requests for notifications.
type NotificationHandler struct {
	service *services.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// GetUserNotifications handles the request to get all notifications for the current user.
func (h *NotificationHandler) GetUserNotifications(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	notifications, err := h.service.GetUserNotifications(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve notifications"})
	}

	return c.JSON(http.StatusOK, notifications)
}

// MarkAsRead handles the request to mark a notification as read.
func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid notification ID"})
	}

	err = h.service.MarkNotificationAsRead(uint(notificationID), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to mark notification as read"})
	}

	return c.NoContent(http.StatusNoContent)
}

// MarkAllAsRead handles the request to mark all of a user's notifications as read.
func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}

	err = h.service.MarkAllUserNotificationsAsRead(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to mark all notifications as read"})
	}

	return c.NoContent(http.StatusNoContent)
}
