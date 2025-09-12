package models

import "time"

type Task struct {
	ID             uint   `gorm:"primaryKey"`
	Title          string `gorm:"not null"`
	Description    string
	UserStoryID    uint      `gorm:"not null"`
	UserStory      UserStory `gorm:"foreignKey:UserStoryID"`
	Status         string    `gorm:"not null;default:'todo'"`
	AssignedToID   *uint
	AssignedTo     *User `gorm:"foreignKey:AssignedToID"`
	EstimatedHours *int
	SpentHours     *int
	IsDeliverable  bool      `gorm:"default:false"`
	CreatedByID    uint      `gorm:"not null"`
	CreatedBy      User      `gorm:"foreignKey:CreatedByID"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	History        []TaskHistory
	Comments       []TaskComment
}

type TaskHistory struct {
	ID          uint   `gorm:"primaryKey"`
	TaskID      uint   `gorm:"not null"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
	Type        string `gorm:"not null"`
	Field       string
	OldValue    string
	NewValue    string
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type TaskComment struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
