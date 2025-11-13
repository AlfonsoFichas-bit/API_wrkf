package models

import "time"

// TaskHistory represents a record of a change made to a task.
type TaskHistory struct {
	ID          uint   `gorm:"primaryKey"`
	TaskID      uint   `gorm:"not null"`
	ChangedBy   User   `gorm:"foreignKey:ChangedByID"`
	ChangedByID uint   `gorm:"not null"`
	FieldName   string `gorm:"not null"` // e.g., "status", "assignedTo"
	OldValue    string
	NewValue    string    `gorm:"not null"`
	ChangedAt   time.Time `gorm:"autoCreateTime"`
}
