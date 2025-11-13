package models

import "time"

// RubricStatus defines the possible statuses for a rubric.
type RubricStatus string

const (
	RubricStatusDraft    RubricStatus = "DRAFT"
	RubricStatusActive   RubricStatus = "ACTIVE"
	RubricStatusArchived RubricStatus = "ARCHIVED"
)

// Rubric represents the main structure for an evaluation rubric.
type Rubric struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name" gorm:"not null"`
	Description string            `json:"description"`
	ProjectID   uint              `json:"projectId" gorm:"not null"`
	Project     Project           `json:"project" gorm:"foreignKey:ProjectID"`
	CreatedByID uint              `json:"createdById" gorm:"not null"`
	CreatedBy   User              `json:"createdBy" gorm:"foreignKey:CreatedByID"`
	Status      RubricStatus      `json:"status" gorm:"type:varchar(20);not null;default:'DRAFT'"`
	IsTemplate  bool              `json:"isTemplate" gorm:"default:false"`
	Criteria    []RubricCriterion `json:"criteria" gorm:"foreignKey:RubricID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// RubricCriterion represents a single criterion within a rubric.
type RubricCriterion struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	RubricID    uint                   `json:"rubricId" gorm:"not null"`
	Title       string                 `json:"title" gorm:"not null"`
	Description string                 `json:"description"`
	MaxPoints   float64                `json:"maxPoints" gorm:"not null"`
	Levels      []RubricCriterionLevel `json:"levels" gorm:"foreignKey:CriterionID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// RubricCriterionLevel represents a performance level within a criterion.
type RubricCriterionLevel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CriterionID uint      `json:"criterionId" gorm:"not null"`
	Score       float64   `json:"score" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
