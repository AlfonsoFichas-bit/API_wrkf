package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// NotificationRepository handles database operations for notifications.
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository.
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create stores a new notification in the database.
func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

// GetByUserID retrieves all notifications for a specific user, ordered by creation time.
func (r *NotificationRepository) GetByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).Order("id ASC").Find(&notifications).Error
	return notifications, err
}

// GetByID retrieves a single notification by its ID.
func (r *NotificationRepository) GetByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, id).Error
	return &notification, err
}

// MarkAsRead marks a single notification as read, ensuring it belongs to the user.
func (r *NotificationRepository) MarkAsRead(notificationID uint, userID uint) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead marks all of a user's notifications as read.
func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// Delete removes a notification, ensuring it belongs to the user.
func (r *NotificationRepository) Delete(notificationID uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&models.Notification{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Or a custom "not found or not authorized" error
	}
	return nil
}
