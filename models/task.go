package models

import "time"

// TaskStatus defines the possible statuses for a Task.
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusInReview   TaskStatus = "in_review"
	StatusDone       TaskStatus = "done"
)

// IsValidTaskStatus checks if a given string is a valid task status.
func IsValidTaskStatus(status string) bool {
	switch TaskStatus(status) {
	case StatusTodo, StatusInProgress, StatusInReview, StatusDone:
		return true
	default:
		return false
	}
}

// Task represents a single task, which is a breakdown of a UserStory.
type Task struct {
	ID             uint   `gorm:"primaryKey"`
	Title          string `gorm:"not null"`
	Description    string
	Priority       string    `gorm:"type:varchar(20);default:'medium'"`
	DueDate        *time.Time
	UserStoryID    uint       `gorm:"not null"`
	UserStory      UserStory  `gorm:"foreignKey:UserStoryID"`
	Status         TaskStatus `gorm:"type:varchar(20);not null;default:'todo'"`
	AssignedToID   *uint
	AssignedTo     *User `gorm:"foreignKey:AssignedToID"`
	EstimatedHours *float32
	SpentHours     *float32
	IsDeliverable  bool      `gorm:"default:false"`

	// Fields for evaluation tracking
	SubmittedForEvaluation bool       `gorm:"default:false"`
	SubmittedAt            *time.Time
	SubmittedByID          *uint
	SubmittedBy            *User `gorm:"foreignKey:SubmittedByID"`

	CreatedByID    uint      `gorm:"not null"`
	CreatedBy      User      `gorm:"foreignKey:CreatedByID"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	History        []TaskHistory
	Comments       []TaskComment
}
