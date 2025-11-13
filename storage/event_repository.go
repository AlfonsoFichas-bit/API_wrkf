package storage

import (
	"time"

	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// EventRepository handles database operations for events.
type EventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new EventRepository.
func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create creates a new event in the database.
func (r *EventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

// FindByID retrieves a single event by its ID.
func (r *EventRepository) FindByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.Preload("CreatedBy").First(&event, id).Error
	return &event, err
}

// FindByProjectAndDateRange retrieves all events for a given project within a specific date range.
func (r *EventRepository) FindByProjectAndDateRange(projectID uint, start time.Time, end time.Time) ([]models.Event, error) {
	var events []models.Event
	err := r.db.
		Where("project_id = ? AND start_date < ? AND end_date > ?", projectID, end, start).
		Preload("CreatedBy").
		Order("start_date ASC").
		Find(&events).Error
	return events, err
}

// Update updates an existing event in the database.
func (r *EventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

// Delete removes an event from the database.
func (r *EventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}
