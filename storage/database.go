package storage

import (
	"fmt"
	"github.com/buga/API_wrkf/config"
	"github.com/buga/API_wrkf/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewConnection creates a new database connection based on the config.
func NewConnection(config *config.DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dsn := config.DSN()

	switch config.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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
		&models.Grade{},
		&models.Notification{},
	)
}
