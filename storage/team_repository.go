package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// TeamRepository handles database operations for teams.
type TeamRepository struct {
	DB *gorm.DB
}

// NewTeamRepository creates a new instance of TeamRepository.
func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{DB: db}
}

// GetTeamsByProjectID retrieves all teams for a given project.
func (r *TeamRepository) GetTeamsByProjectID(projectID uint) ([]models.Team, error) {
	var teams []models.Team
	err := r.DB.Where("project_id = ?", projectID).Preload("Members").Find(&teams).Error
	return teams, err
}
