package storage

import (
	"time"
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// DeadlineRepository handles database operations for deadlines.
type DeadlineRepository struct {
	DB *gorm.DB
}

// NewDeadlineRepository creates a new instance of DeadlineRepository.
func NewDeadlineRepository(db *gorm.DB) *DeadlineRepository {
	return &DeadlineRepository{DB: db}
}

// GetUpcomingDeadlines retrieves deadlines within a given timeframe, filtered by projects and type.
func (r *DeadlineRepository) GetUpcomingDeadlines(from time.Time, to time.Time, projectIDs []uint, deadlineType *string) ([]models.Deadline, error) {
	var deadlines []models.Deadline

	query := r.DB.
		Where("is_active = ?", true).
		Where("date >= ?", from).
		Where("date <= ?", to)

	if len(projectIDs) > 0 {
		// Include deadlines for the user's projects AND global deadlines (project_id IS NULL).
		query = query.Where("project_id IN ? OR project_id IS NULL", projectIDs)
	} else {
		// If user has no projects, only show global deadlines.
		query = query.Where("project_id IS NULL")
	}

	if deadlineType != nil && *deadlineType != "all" && *deadlineType != "" {
		query = query.Where("type = ?", *deadlineType)
	}

	err := query.Order("date ASC").Find(&deadlines).Error
	return deadlines, err
}
