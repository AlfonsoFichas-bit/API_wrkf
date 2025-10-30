package models

import "time"

// Attachment represents a file attached to a project, task, etc.
type Attachment struct {
	ID        uint      `gorm:"primaryKey"`
	FileName  string    `gorm:"not null"`
	FilePath  string    `gorm:"not null"`
	FileType  string
	FileSize  int64
	ProjectID *uint // Optional: link to a project
	TaskID    *uint // Optional: link to a task
	UploadedByID uint
	UploadedBy   User      `gorm:"foreignKey:UploadedByID"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
