package models

import "time"

// TaskComment represents a comment made on a task.
type TaskComment struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	AuthorID  uint      `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
