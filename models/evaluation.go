package models

import "time"

type Rubric struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	ProjectID   uint              `gorm:"not null"`
	Project     Project           `gorm:"foreignKey:ProjectID"`
	CreatedByID uint              `gorm:"not null"`
	CreatedBy   User              `gorm:"foreignKey:CreatedByID"`
	IsTemplate  bool              `gorm:"default:false"`
	Status      string            `gorm:"not null;default:'draft'"`
	CreatedAt   time.Time         `gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime"`
	Criteria    []RubricCriterium `gorm:"foreignKey:RubricID"`
}

type RubricCriterium struct {
	ID          uint   `gorm:"primaryKey"`
	RubricID    uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	MaxPoints   int                    `gorm:"not null"`
	CreatedAt   time.Time              `gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `gorm:"autoUpdateTime"`
	Levels      []RubricCriterionLevel `gorm:"foreignKey:CriterionID"`
}

type RubricCriterionLevel struct {
	ID          uint      `gorm:"primaryKey"`
	CriterionID uint      `gorm:"not null"`
	Description string    `gorm:"not null"`
	PointValue  int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

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
	Criterium    RubricCriterium `gorm:"foreignKey:CriterionID"`
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
