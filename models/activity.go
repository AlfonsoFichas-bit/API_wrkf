package models

import "time"

// Activity represents a recent action within a project for the dashboard feed.
type Activity struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type"` // e.g., 'task_completed', 'comment', 'sprint_started'
	UserID      uint      `json:"userId"`
	User        User      `gorm:"foreignKey:UserID"`
	EntityType  string    `gorm:"type:varchar(50);not null" json:"entityType"` // 'task', 'user_story', 'sprint'
	EntityID    uint      `gorm:"not null" json:"entityId"`
	ProjectID   uint      `json:"projectId"`
	Project     Project   `gorm:"foreignKey:ProjectID"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"timestamp"`
}
