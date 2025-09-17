package models

import "time"

type Sprint struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Goal        string
	ProjectID   uint    `gorm:"not null"`
	Project     Project `gorm:"foreignKey:ProjectID"`
	Status      string  `gorm:"not null;default:'planned'"`
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedByID uint        `gorm:"not null"`
	CreatedBy   User        `gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`
	UserStories []UserStory `gorm:"foreignKey:SprintID"` // <-- RELATIONSHIP ADDED
}
