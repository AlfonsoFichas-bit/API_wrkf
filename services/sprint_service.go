package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

type SprintService struct {
	sprintRepo    *storage.SprintRepository
	userStoryRepo *storage.UserStoryRepository
}

// NewSprintService creates a new instance of SprintService.
func NewSprintService(sprintRepo *storage.SprintRepository, userStoryRepo *storage.UserStoryRepository) *SprintService {
	return &SprintService{sprintRepo: sprintRepo, userStoryRepo: userStoryRepo}
}

// CreateSprint handles the business logic for creating a new sprint.
func (s *SprintService) CreateSprint(sprint *models.Sprint, projectID uint, creatorID uint) error {
	sprint.ProjectID = projectID
	sprint.CreatedByID = creatorID
	return s.sprintRepo.CreateSprint(sprint)
}

// GetSprintsByProjectID retrieves all sprints for a specific project.
func (s *SprintService) GetSprintsByProjectID(projectID uint) ([]models.Sprint, error) {
	return s.sprintRepo.GetSprintsByProjectID(projectID)
}

// GetSprintByID retrieves a single sprint by its ID.
func (s *SprintService) GetSprintByID(id uint) (*models.Sprint, error) {
	return s.sprintRepo.GetSprintByID(id)
}

// UpdateSprint handles the business logic for updating a sprint.
func (s *SprintService) UpdateSprint(sprint *models.Sprint) error {
	// In the future, you could add permission checks here.
	return s.sprintRepo.UpdateSprint(sprint)
}

func (s *SprintService) DeleteSprint(id uint) error {
	return s.sprintRepo.DeleteSprint(id)
}

// AddUserStoryToSprint assigns a user story to a sprint.
func (s *SprintService) AddUserStoryToSprint(sprintID uint, userStoryID uint) error {
	return s.userStoryRepo.AssignUserStoryToSprint(userStoryID, sprintID)
}
