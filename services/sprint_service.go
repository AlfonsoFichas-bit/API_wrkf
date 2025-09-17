package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// SprintService handles the business logic for sprints.
type SprintService struct {
	Repo *storage.SprintRepository
}

// NewSprintService creates a new instance of SprintService.
func NewSprintService(repo *storage.SprintRepository) *SprintService {
	return &SprintService{Repo: repo}
}

// CreateSprint handles the business logic for creating a new sprint.
func (s *SprintService) CreateSprint(sprint *models.Sprint, projectID uint, creatorID uint) error {
	sprint.ProjectID = projectID
	sprint.CreatedByID = creatorID
	return s.Repo.CreateSprint(sprint)
}

// GetSprintsByProjectID retrieves all sprints for a specific project.
func (s *SprintService) GetSprintsByProjectID(projectID uint) ([]models.Sprint, error) {
	return s.Repo.GetSprintsByProjectID(projectID)
}

// GetSprintByID retrieves a single sprint by its ID.
func (s *SprintService) GetSprintByID(id uint) (*models.Sprint, error) {
	return s.Repo.GetSprintByID(id)
}

// UpdateSprint handles the business logic for updating a sprint.
func (s *SprintService) UpdateSprint(sprint *models.Sprint) error {
	// In the future, you could add permission checks here.
	return s.Repo.UpdateSprint(sprint)
}

// DeleteSprint handles the business logic for deleting a sprint.
func (s *SprintService) DeleteSprint(id uint) error {
	// In the future, you could add logic here to move user stories back to the backlog.
	return s.Repo.DeleteSprint(id)
}
