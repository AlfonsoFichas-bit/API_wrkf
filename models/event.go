package models

import "time"

// Event represents a calendar event, such as a meeting, deadline, or milestone.
type Event struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate" gorm:"not null"`
	EndDate     time.Time `json:"endDate" gorm:"not null"`
	Type        string    `json:"type" gorm:"default:'general'"` // e.g., 'meeting', 'deadline', 'milestone'
	ProjectID   uint      `json:"projectId" gorm:"not null"`
	Project     Project   `json:"project" gorm:"foreignKey:ProjectID"`
	CreatedByID uint      `json:"createdById" gorm:"not null"`
	CreatedBy   User      `json:"createdBy" gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
