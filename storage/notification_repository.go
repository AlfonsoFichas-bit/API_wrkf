package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// NotificationRepository handles database operations for notifications.
type NotificationRepository struct {
	DB *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository.
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// Create creates a new notification in the database.
func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.DB.Create(notification).Error
}

// GetByUserID retrieves all notifications for a specific user.
func (r *NotificationRepository) GetByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

// MarkAsRead marks a single notification as read.
func (r *NotificationRepository) MarkAsRead(notificationID uint, userID uint) error {
	return r.DB.Model(&models.Notification{}).Where("id = ? AND user_id = ?", notificationID, userID).Update("is_read", true).Error
}

// MarkAllAsRead marks all unread notifications for a user as read.
func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return r.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}
