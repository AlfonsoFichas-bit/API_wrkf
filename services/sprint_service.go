package services

import (
	"API_wrkf/models"
	"API_wrkf/storage"
)

type SprintService struct {
	Repo *storage.SprintRepository
}

func NewSprintService(repo *storage.SprintRepository) *SprintService {
	return &SprintService{Repo: repo}
}

func (s *SprintService) CreateSprint(sprint *models.Sprint, userID uint) error {
	sprint.CreatedByID = userID
	return s.Repo.CreateSprint(sprint)
}
