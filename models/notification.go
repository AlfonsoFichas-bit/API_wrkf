package models

import (
	"gorm.io/gorm"
)

// Notification represents a notification for a user.
type Notification struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index" json:"user_id"` // ID of the user receiving the notification
	Message string `gorm:"not null" json:"message"`       // The content of the notification
	IsRead  bool   `gorm:"default:false" json:"is_read"`  // Whether the notification has been read
	Link    string `json:"link,omitempty"`                // A URL to the relevant resource, e.g., /projects/123
}
