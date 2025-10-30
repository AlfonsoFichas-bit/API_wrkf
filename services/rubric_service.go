package services

import (
	"fmt"
	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"
)

// RubricService defines the business logic for rubrics.
type RubricService interface {
	CreateRubric(rubric *models.Rubric) error
	GetAllRubrics(filters map[string]interface{}) ([]models.Rubric, error)
	GetRubricByID(id uint) (*models.Rubric, error)
	UpdateRubric(rubric *models.Rubric) error
	DeleteRubric(id uint) error
	DuplicateRubric(id uint) (*models.Rubric, error)
}

type rubricService struct {
	repo storage.IRubricRepository
}

// NewRubricService creates a new instance of RubricService.
func NewRubricService(repo storage.IRubricRepository) RubricService {
	return &rubricService{repo: repo}
}

func (s *rubricService) CreateRubric(rubric *models.Rubric) error {
	return s.repo.CreateRubric(rubric)
}

func (s *rubricService) GetAllRubrics(filters map[string]interface{}) ([]models.Rubric, error) {
	return s.repo.GetAllRubrics(filters)
}

func (s *rubricService) GetRubricByID(id uint) (*models.Rubric, error) {
	return s.repo.GetRubricByID(id)
}

func (s *rubricService) UpdateRubric(rubric *models.Rubric) error {
	return s.repo.UpdateRubric(rubric)
}

func (s *rubricService) DeleteRubric(id uint) error {
	return s.repo.DeleteRubric(id)
}

func (s *rubricService) DuplicateRubric(id uint) (*models.Rubric, error) {
	original, err := s.repo.GetRubricByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find original rubric: %w", err)
	}

	// Create a deep copy for duplication
	newRubric := *original
	newRubric.ID = 0 // Set ID to 0 to create a new record
	newRubric.Name = "Copia de " + original.Name
	newRubric.Status = models.RubricStatusDraft // New duplicates are always drafts

	// Deep copy criteria and levels
	newRubric.Criteria = make([]models.RubricCriterion, len(original.Criteria))
	for i, crit := range original.Criteria {
		newCrit := crit
		newCrit.ID = 0
		newCrit.RubricID = 0 // Will be set by GORM on creation

		newCrit.Levels = make([]models.RubricCriterionLevel, len(crit.Levels))
		for j, level := range crit.Levels {
			newLevel := level
			newLevel.ID = 0
			newLevel.CriterionID = 0 // Will be set by GORM on creation
			newCrit.Levels[j] = newLevel
		}
		newRubric.Criteria[i] = newCrit
	}

	if err := s.repo.CreateRubric(&newRubric); err != nil {
		return nil, fmt.Errorf("failed to create duplicated rubric: %w", err)
	}

	return &newRubric, nil
}
