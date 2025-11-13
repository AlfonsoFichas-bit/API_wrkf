package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// NotificationService provides notification-related services.
type NotificationService struct {
	repo *storage.NotificationRepository
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(repo *storage.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// CreateNotification creates and stores a new notification.
func (s *NotificationService) CreateNotification(userID uint, message string, link string) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  userID,
		Message: message,
		Link:    link,
		IsRead:  false,
	}

	err := s.repo.Create(notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

// GetUserNotifications retrieves all notifications for a specific user.
func (s *NotificationService) GetUserNotifications(userID uint) ([]models.Notification, error) {
	return s.repo.GetByUserID(userID)
}

// MarkNotificationAsRead marks a single notification as read.
func (s *NotificationService) MarkNotificationAsRead(notificationID uint, userID uint) error {
	return s.repo.MarkAsRead(notificationID, userID)
}

// MarkAllUserNotificationsAsRead marks all of a user's notifications as read.
func (s *NotificationService) MarkAllUserNotificationsAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

// GetNotificationByID retrieves a single notification by its ID.
func (s *NotificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	return s.repo.GetByID(id)
}

// DeleteNotification deletes a notification.
func (s *NotificationService) DeleteNotification(notificationID uint, userID uint) error {
	return s.repo.Delete(notificationID, userID)
}
