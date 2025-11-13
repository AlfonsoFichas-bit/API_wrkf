package services

import (
	"fmt"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// EventService handles the business logic for calendar events.
type EventService struct {
	EventRepo      *storage.EventRepository
	ProjectService *ProjectService // To check user roles
}

// NewEventService creates a new instance of EventService.
func NewEventService(eventRepo *storage.EventRepository, projectService *ProjectService) *EventService {
	return &EventService{
		EventRepo:      eventRepo,
		ProjectService: projectService,
	}
}

// checkUserPermission is a helper to verify if a user is a member of the project.
func (s *EventService) checkUserPermission(userID, projectID uint) error {
	_, err := s.ProjectService.GetUserRoleInProject(userID, projectID)
	if err != nil {
		return fmt.Errorf("user does not have permission in this project")
	}
	return nil
}

// CreateEvent handles the business logic for creating a new event.
func (s *EventService) CreateEvent(event *models.Event, projectID, creatorID uint) (*models.Event, error) {
	if err := s.checkUserPermission(creatorID, projectID); err != nil {
		return nil, err
	}
	if event.StartDate.After(event.EndDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	event.ProjectID = projectID
	event.CreatedByID = creatorID

	if err := s.EventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}
	return s.EventRepo.FindByID(event.ID) // Return hydrated event
}

// GetEventByID retrieves a single event, checking for permissions.
func (s *EventService) GetEventByID(eventID, userID uint) (*models.Event, error) {
	event, err := s.EventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}
	if err := s.checkUserPermission(userID, event.ProjectID); err != nil {
		return nil, err
	}
	return event, nil
}

// GetEventsForProject retrieves events for a project within a date range, checking permissions.
func (s *EventService) GetEventsForProject(projectID, userID uint, start, end time.Time) ([]models.Event, error) {
	if err := s.checkUserPermission(userID, projectID); err != nil {
		return nil, err
	}
	return s.EventRepo.FindByProjectAndDateRange(projectID, start, end)
}

// UpdateEvent handles updating an event, checking for permissions.
func (s *EventService) UpdateEvent(eventID, userID uint, updates map[string]interface{}) (*models.Event, error) {
	event, err := s.EventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found")
	}
	if err := s.checkUserPermission(userID, event.ProjectID); err != nil {
		return nil, err
	}

	// Apply updates
	if title, ok := updates["title"].(string); ok {
		event.Title = title
	}
	if description, ok := updates["description"].(string); ok {
		event.Description = description
	}
	// Add other fields as needed, e.g., dates

	if err := s.EventRepo.Update(event); err != nil {
		return nil, fmt.Errorf("could not update event: %w", err)
	}
	return event, nil
}

// DeleteEvent handles deleting an event, checking for permissions.
func (s *EventService) DeleteEvent(eventID, userID uint) error {
	event, err := s.EventRepo.FindByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if err := s.checkUserPermission(userID, event.ProjectID); err != nil {
		return err
	}
	return s.EventRepo.Delete(eventID)
}
