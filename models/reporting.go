package models

import (
	"encoding/json"
	"time"
)

type Report struct {
	ID                   uint   `gorm:"primaryKey"`
	Title                string `gorm:"not null"`
	Description          string
	Type                 string `gorm:"not null"`
	ProjectID            *uint
	Project              *Project `gorm:"foreignKey:ProjectID"`
	SprintID             *uint
	Sprint               *Sprint `gorm:"foreignKey:SprintID"`
	UserID               *uint
	User                 *User `gorm:"foreignKey:UserID"`
	StartDate            *time.Time
	EndDate              *time.Time
	IncludeBurndown      *bool
	IncludeVelocity      *bool
	IncludeUserMetrics   *bool
	IncludeProjectHealth *bool
	CustomSections       []string        `gorm:"type:text[]"`
	CreatedByID          uint            `gorm:"not null"`
	CreatedBy            User            `gorm:"foreignKey:CreatedByID"`
	GeneratedAt          time.Time       `gorm:"autoCreateTime"`
	Data                 json.RawMessage `gorm:"type:jsonb"`
	ExportFormats        []string        `gorm:"type:text[]"`
	CreatedAt            time.Time       `gorm:"autoCreateTime"`
	UpdatedAt            time.Time       `gorm:"autoUpdateTime"`
}

type ScheduledReport struct {
	ID             uint   `gorm:"primaryKey"`
	ReportConfigID uint   `gorm:"not null"`
	ReportConfig   Report `gorm:"foreignKey:ReportConfigID"`
	Frequency      string `gorm:"not null"`
	NextRunTime    *time.Time
	CreatedByID    uint     `gorm:"not null"`
	CreatedBy      User     `gorm:"foreignKey:CreatedByID"`
	Recipients     []string `gorm:"type:text[]"`
	LastRunTime    *time.Time
	LastReportID   *uint
	LastReport     *Report   `gorm:"foreignKey:LastReportID"`
	ExportFormats  []string  `gorm:"type:text[]"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type ProjectMetric struct {
	ID                   uint      `gorm:"primaryKey"`
	ProjectID            uint      `gorm:"not null"`
	Project              Project   `gorm:"foreignKey:ProjectID"`
	Date                 time.Time `gorm:"not null"`
	TotalUserStories     *int
	CompletedUserStories *int
	TotalPoints          *int
	CompletedPoints      *int
	AverageVelocity      *int
	PredictedCompletion  *time.Time
	HealthScore          *int
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
}

type SprintMetric struct {
	ID              uint      `gorm:"primaryKey"`
	SprintID        uint      `gorm:"not null"`
	Sprint          Sprint    `gorm:"foreignKey:SprintID"`
	Date            time.Time `gorm:"not null"`
	TotalPoints     *int
	CompletedPoints *int
	RemainingPoints *int
	TasksCompleted  *int
	TasksRemaining  *int
	IdealBurndown   *int
	ProjectID       *uint
	Project         *Project  `gorm:"foreignKey:ProjectID"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

type UserMetric struct {
	ID                uint `gorm:"primaryKey"`
	UserID            uint `gorm:"not null"`
	User              User `gorm:"foreignKey:UserID"`
	SprintID          *uint
	Sprint            *Sprint   `gorm:"foreignKey:SprintID"`
	Date              time.Time `gorm:"not null"`
	TasksCompleted    *int
	PointsContributed *int
	HoursLogged       *int
	Efficiency        *int
	ProjectID         *uint
	Project           *Project  `gorm:"foreignKey:ProjectID"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}
