package models

import "time"

// Report represents a generated report.
type Report struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Type      string    // e.g., "burndown", "velocity"
	ProjectID uint      `gorm:"not null"`
	Data      string    `gorm:"type:jsonb"` // Store report data as JSON
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// ScheduledReport defines a report that is generated automatically.
type ScheduledReport struct {
	ID         uint   `gorm:"primaryKey"`
	ReportType string `gorm:"not null"`
	Schedule   string `gorm:"not null"` // e.g., "weekly", "monthly"
	ProjectID  uint   `gorm:"not null"`
}

// ProjectMetric stores metrics for a project.
type ProjectMetric struct {
	ID        uint `gorm:"primaryKey"`
	ProjectID uint `gorm:"not null;uniqueIndex"`
	// Add metric fields here, e.g., TotalTasks, CompletedTasks, etc.
}

// SprintMetric stores metrics for a sprint.
type SprintMetric struct {
	ID       uint `gorm:"primaryKey"`
	SprintID uint `gorm:"not null;uniqueIndex"`
	// Add metric fields here, e.g., Velocity, Burndown, etc.
}

// UserMetric stores metrics for a user.
type UserMetric struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null;uniqueIndex"`
	// Add metric fields here, e.g., TasksCompleted, AverageEffort, etc.
}
