package models

import "time"

// Deadline represents an upcoming deadline for the dashboard.
type Deadline struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	Type       string    `gorm:"type:varchar(50);not null" json:"type"` // 'sprint_end', 'evaluation_deadline', 'project_milestone'
	ProjectID  *uint     `json:"projectId"`
	Project    *Project  `gorm:"foreignKey:ProjectID"`
	EntityType *string   `gorm:"type:varchar(50)" json:"entityType,omitempty"` // 'sprint', 'evaluation'
	EntityID   *uint     `json:"entityId,omitempty"`
	Date       time.Time `gorm:"not null" json:"date"`
	IsActive   bool      `gorm:"default:true" json:"-"` // To hide or show deadlines
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"-"`
}
