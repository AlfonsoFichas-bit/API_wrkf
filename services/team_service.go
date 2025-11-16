package services

import (
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// TeamService handles the business logic for teams.
type TeamService struct {
	Repo *storage.TeamRepository
}

// NewTeamService creates a new instance of TeamService.
func NewTeamService(repo *storage.TeamRepository) *TeamService {
	return &TeamService{Repo: repo}
}

// GetTeamsByProjectID retrieves all teams for a given project.
func (s *TeamService) GetTeamsByProjectID(projectID uint) ([]models.Team, error) {
	return s.Repo.GetTeamsByProjectID(projectID)
}
