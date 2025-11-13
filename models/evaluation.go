package models

import "time"

// Evaluation represents a formal assessment of a submitted task (deliverable)
// against a specific rubric.
type Evaluation struct {
	ID                   uint   `gorm:"primaryKey"`
	TaskID               uint   `gorm:"not null;uniqueIndex:idx_task_evaluator"`
	Task                 Task   `gorm:"foreignKey:TaskID"`
	EvaluatorID          uint   `gorm:"not null;uniqueIndex:idx_task_evaluator"`
	Evaluator            User   `gorm:"foreignKey:EvaluatorID"`
	RubricID             uint   `gorm:"not null"`
	Rubric               Rubric `gorm:"foreignKey:RubricID"`
	OverallFeedback      string `gorm:"type:text"`
	TotalScore           float64
	Status               string                `gorm:"type:varchar(20);not null;default:'draft'"` // e.g., draft, published
	CreatedAt            time.Time             `gorm:"autoCreateTime"`
	UpdatedAt            time.Time             `gorm:"autoUpdateTime"`
	CriterionEvaluations []CriterionEvaluation `gorm:"foreignKey:EvaluationID;constraint:OnDelete:CASCADE;"`
}

// CriterionEvaluation stores the score and feedback for a single criterion
// within a larger evaluation.
type CriterionEvaluation struct {
	ID           uint            `gorm:"primaryKey"`
	EvaluationID uint            `gorm:"not null"`
	CriterionID  uint            `gorm:"not null"`
	Criterion    RubricCriterion `gorm:"foreignKey:CriterionID"`
	Score        float64         `gorm:"not null"`
	Feedback     string          `gorm:"type:text"`
	CreatedAt    time.Time       `gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime"`
}
