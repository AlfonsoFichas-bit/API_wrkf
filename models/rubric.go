package models

import "gorm.io/gorm"

// RubricStatus defines the status of a rubric.
type RubricStatus string

const (
    RubricStatusDraft     RubricStatus = "draft"
    RubricStatusPublished RubricStatus = "published"
)

// Rubric represents the rubric for an evaluation
type Rubric struct {
    gorm.Model
    Name        string             `json:"name"`
    Description string             `json:"description"`
    ProjectID   uint               `json:"project_id"`
    Status      RubricStatus       `json:"status"`
    Criteria    []RubricCriterion  `json:"criteria" gorm:"foreignKey:RubricID"`
}

// RubricCriterion represents a single criterion in a rubric.
type RubricCriterion struct {
    gorm.Model
    Description string                 `json:"description"`
    RubricID    uint                   `json:"rubric_id"`
    Levels      []RubricCriterionLevel `json:"levels" gorm:"foreignKey:CriterionID"`
}

// RubricCriterionLevel represents a performance level for a criterion.
type RubricCriterionLevel struct {
    gorm.Model
    Description string `json:"description"`
    Points      uint   `json:"points"`
    CriterionID uint   `json:"criterion_id"`
}
