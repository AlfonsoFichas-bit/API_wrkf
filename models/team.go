package models

import "time"

// Team represents a team of users within a project.
type Team struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	ProjectID uint      `gorm:"not null" json:"projectId"`
	Project   Project   `gorm:"foreignKey:ProjectID"`
	Members   []User    `gorm:"many2many:team_members;" json:"members"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
