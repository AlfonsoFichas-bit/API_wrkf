package models

import "time"

type Project struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Status      string `gorm:"not null;default:'planning'"`
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedByID uint      `gorm:"not null"`
	CreatedBy   User      `gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Members     []ProjectMember
}

type ProjectMember struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	ProjectID uint      `gorm:"not null"`
	Role      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
