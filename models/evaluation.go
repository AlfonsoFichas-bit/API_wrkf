package models

import "time"

type Evaluation struct {
	ID                   uint   `gorm:"primaryKey"`
	DeliverableID        uint   `gorm:"not null"`
	Deliverable          Task   `gorm:"foreignKey:DeliverableID"`
	EvaluatorID          uint   `gorm:"not null"`
	Evaluator            User   `gorm:"foreignKey:EvaluatorID"`
	StudentID            uint   `gorm:"not null"`
	Student              User   `gorm:"foreignKey:StudentID"`
	RubricID             uint   `gorm:"not null"`
	Rubric               Rubric `gorm:"foreignKey:RubricID"`
	OverallFeedback      string
	TotalScore           int    `gorm:"not null"`
	MaxPossibleScore     int    `gorm:"not null"`
	Status               string `gorm:"not null;default:'draft'"`
	EvaluatedAt          *time.Time
	CreatedAt            time.Time             `gorm:"autoCreateTime"`
	UpdatedAt            time.Time             `gorm:"autoUpdateTime"`
	CriterionEvaluations []CriterionEvaluation `gorm:"foreignKey:EvaluationID"`
}

type CriterionEvaluation struct {
	ID           uint            `gorm:"primaryKey"`
	EvaluationID uint            `gorm:"not null"`
	CriterionID  uint            `gorm:"not null"`
	Criterium    RubricCriterion `gorm:"foreignKey:CriterionID"` // Changed from RubricCriterium
	Score        int             `gorm:"not null"`
	Feedback     string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type Attachment struct {
	ID            uint      `gorm:"primaryKey"`
	DeliverableID uint      `gorm:"not null"`
	Deliverable   Task      `gorm:"foreignKey:DeliverableID"`
	FileName      string    `gorm:"not null"`
	FileType      string    `gorm:"not null"`
	FileSize      int       `gorm:"not null"`
	URL           string    `gorm:"not null"`
	UploadedByID  uint      `gorm:"not null"`
	UploadedBy    User      `gorm:"foreignKey:UploadedByID"`
	UploadedAt    time.Time `gorm:"autoCreateTime"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}