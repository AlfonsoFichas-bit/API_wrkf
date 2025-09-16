package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

type SprintRepository struct {
	DB *gorm.DB
}

func NewSprintRepository(db *gorm.DB) *SprintRepository {
	return &SprintRepository{DB: db}
}

func (r *SprintRepository) CreateSprint(sprint *models.Sprint) error {
	return r.DB.Create(sprint).Error
}
