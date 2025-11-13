package models

import "time"

type UserStory struct {
	ID                 uint   `gorm:"primaryKey"`
	Title              string `gorm:"not null"`
	Description        string `gorm:"not null"`
	AcceptanceCriteria string `gorm:"not null"`
	Priority           string `gorm:"not null;default:'medium'"`
	Status             string `gorm:"not null;default:'backlog'"`
	Points             *int
	ProjectID          uint    `gorm:"not null"`
	Project            Project `gorm:"foreignKey:ProjectID"`
	SprintID           *uint   // Pointer to allow null values
	Sprint             *Sprint // Let GORM infer the relationship via convention
	CreatedByID        uint    `gorm:"not null"`
	CreatedBy          User    `gorm:"foreignKey:CreatedByID"`
	AssignedToID       *uint
	AssignedTo         *User     `gorm:"foreignKey:AssignedToID"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}
