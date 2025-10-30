package models

import (
	"gorm.io/gorm"
)

// Evaluation represents a student's evaluation for a specific task or deliverable.
type Evaluation struct {
	gorm.Model
	StudentID   uint   `json:"student_id"`
	ProjectID   uint   `json:"project_id"`
	TaskID      uint   `json:"task_id"`
	RubricID    uint   `json:"rubric_id"`
	FinalGrade  float64 `json:"final_grade"`
	Comments    string `json:"comments"`
	Grades      []Grade `json:"grades" gorm:"foreignKey:EvaluationID"`
}

// Grade represents the score for a specific criterion in an evaluation.
type Grade struct {
	gorm.Model
	EvaluationID uint `json:"evaluation_id"`
	CriterionID  uint `json:"criterion_id"`
	Score        uint `json:"score"`
}
