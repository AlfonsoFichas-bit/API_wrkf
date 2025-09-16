
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
func (s *SprintService) CreateSprint(sprint *models.Sprint, creatorID uint) error {
	sprint.CreatedByID = creatorID
	return s.Repo.CreateSprint(sprint)
}
