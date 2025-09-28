package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// BacklogService handles business logic for the backlog.
type BacklogService struct {
	userStoryRepo *storage.UserStoryRepository
}

// NewBacklogService creates a new instance of BacklogService.
func NewBacklogService(userStoryRepo *storage.UserStoryRepository) *BacklogService {
	return &BacklogService{userStoryRepo: userStoryRepo}
}

// GetProductBacklog retrieves all user stories for a given project that are not assigned to a sprint.
func (s *BacklogService) GetProductBacklog(projectID uint) ([]models.UserStory, error) {
	return s.userStoryRepo.GetBacklogUserStoriesByProjectID(projectID)
}

// UpdateUserStoryStatus updates the status of a user story.
func (s *BacklogService) UpdateUserStoryStatus(id uint, status string) error {
	return s.userStoryRepo.UpdateUserStoryStatus(id, status)
}