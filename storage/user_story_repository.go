package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

// UserStoryRepository handles database operations for user stories.
type UserStoryRepository struct {
	DB *gorm.DB
}

// NewUserStoryRepository creates a new instance of UserStoryRepository.
func NewUserStoryRepository(db *gorm.DB) *UserStoryRepository {
	return &UserStoryRepository{DB: db}
}

// CreateUserStory adds a new user story to the database.
func (r *UserStoryRepository) CreateUserStory(userStory *models.UserStory) error {
	return r.DB.Create(userStory).Error
}

// GetUserStoriesByProjectID retrieves all user stories for a given project ID.
func (r *UserStoryRepository) GetUserStoriesByProjectID(projectID uint) ([]models.UserStory, error) {
	var userStories []models.UserStory
	err := r.DB.Where("project_id = ?", projectID).Preload("CreatedBy").Find(&userStories).Error
	return userStories, err
}

// GetUserStoryByID retrieves a single user story by its ID, preloading most related data.
// Sprint is loaded manually in the service layer due to preload issues.
func (r *UserStoryRepository) GetUserStoryByID(id uint) (*models.UserStory, error) {
	var userStory models.UserStory
	err := r.DB.Preload("Project").Preload("CreatedBy").Preload("AssignedTo").First(&userStory, id).Error
	return &userStory, err
}

// UpdateUserStory updates an existing user story in the database.
func (r *UserStoryRepository) UpdateUserStory(userStory *models.UserStory) error {
	return r.DB.Save(userStory).Error
}

// DeleteUserStory removes a user story from the database by its ID.
func (r *UserStoryRepository) DeleteUserStory(id uint) error {
	return r.DB.Delete(&models.UserStory{}, id).Error
}
