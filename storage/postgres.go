package storage

import (
	"fmt"

	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // Importar el driver SQLite
	"gorm.io/gorm"
)

// NewConnection creates a new database connection.
func NewConnection(config *config.DBConfig) (*gorm.DB, error) {
	dsn := config.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// NewTestConnection creates a new in-memory SQLite database connection for testing.
func NewTestConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}
	return db, nil
}

// Migrate automates the database migration for all models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Sprint{},
		&models.UserStory{},
		&models.Task{},
		&models.TaskHistory{},
		&models.TaskComment{},
		&models.Rubric{},
		&models.RubricCriterion{},
		&models.RubricCriterionLevel{},
		&models.Evaluation{},
		&models.CriterionEvaluation{},
		&models.Conversation{},
		&models.ConversationMember{},
		&models.Message{},
		&models.MessageAttachment{},
		&models.MessageReadBy{},
		&models.Report{},
		&models.ScheduledReport{},
		&models.ProjectMetric{},
		&models.SprintMetric{},
		&models.UserMetric{},
		&models.Notification{},
		&models.Event{},
	)
}
