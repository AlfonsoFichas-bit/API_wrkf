package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
	"time"
)

// ActivityRepository handles database operations for activities.
type ActivityRepository struct {
	DB *gorm.DB
}

// NewActivityRepository creates a new instance of ActivityRepository.
func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{DB: db}
}

// CreateActivity adds a new activity to the database.
func (r *ActivityRepository) CreateActivity(activity *models.Activity) error {
	return r.DB.Create(activity).Error
}

// GetRecentActivities retrieves recent activities with filters.
// It fetches activities from projects the given user is a member of.
func (r *ActivityRepository) GetRecentActivities(projectIDs []uint, userID *uint, limit int) ([]models.Activity, error) {
	var activities []models.Activity
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	query := r.DB.
		Preload("User").
		Preload("Project").
		Where("project_id IN ?", projectIDs).
		Where("created_at >= ?", sevenDaysAgo).
		Order("created_at DESC").
		Limit(limit)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	err := query.Find(&activities).Error
	return activities, err
}
