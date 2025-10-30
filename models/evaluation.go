package models

import "time"

// Evaluation represents the assessment of a deliverable task.
type Evaluation struct {
	ID          uint      `gorm:"primaryKey"`
	TaskID      uint      `gorm:"uniqueIndex;not null"` // An evaluation is unique to a task
	Task        Task      `gorm:"foreignKey:TaskID"`
	EvaluatorID uint      `gorm:"not null"` // User ID of the evaluator (e.g., Product Owner)
	Evaluator   User      `gorm:"foreignKey:EvaluatorID"`
	RubricID    uint      // Optional: Rubric used for this evaluation
	Rubric      Rubric    `gorm:"foreignKey:RubricID"`
	Score       float32   // Final calculated score
	Comments    string    `gorm:"type:text"` // General feedback or summary
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Criteria    []CriterionEvaluation
	Feedback    []FeedbackComment
}

// CriterionEvaluation stores the score for a single criterion from the rubric.
type CriterionEvaluation struct {
	ID                uint            `gorm:"primaryKey"`
	EvaluationID      uint            `gorm:"not null"`
	RubricCriterionID uint            `gorm:"not null"` // The criterion from the original rubric
	RubricCriterion   RubricCriterion `gorm:"foreignKey:RubricCriterionID"`
	Score             float32         `gorm:"not null"`
	Comment           string          // Optional comment specific to this criterion
}

// FeedbackComment allows for specific, threaded-style comments on an evaluation.
type FeedbackComment struct {
	ID           uint      `gorm:"primaryKey"`
	EvaluationID uint      `gorm:"not null"`
	AuthorID     uint      `gorm:"not null"`
	Author       User      `gorm:"foreignKey:AuthorID"`
	Content      string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
