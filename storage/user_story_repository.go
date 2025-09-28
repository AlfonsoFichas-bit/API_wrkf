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

// GetBacklogUserStoriesByProjectID retrieves all user stories for a given project that are not assigned to a sprint.
func (r *UserStoryRepository) GetBacklogUserStoriesByProjectID(projectID uint) ([]models.UserStory, error) {
	var userStories []models.UserStory
	err := r.DB.Where("project_id = ? AND sprint_id IS NULL", projectID).Preload("CreatedBy").Find(&userStories).Error
	return userStories, err
}

// GetUserStoryByID retrieves a single user story by its ID, preloading all related data.
func (r *UserStoryRepository) GetUserStoryByID(id uint) (*models.UserStory, error) {
	var userStory models.UserStory
	err := r.DB.Preload("Project").Preload("CreatedBy").Preload("AssignedTo").First(&userStory, id).Error
	return &userStory, err
}

// UpdateUserStory updates an existing user story in the database.
func (r *UserStoryRepository) UpdateUserStory(userStory *models.UserStory) error {
	return r.DB.Save(userStory).Error
}

// UpdateUserStoryStatus updates the status of a user story.
func (r *UserStoryRepository) UpdateUserStoryStatus(id uint, status string) error {
	return r.DB.Model(&models.UserStory{}).Where("id = ?", id).Update("status", status).Error
}

// AssignUserStoryToSprint assigns a user story to a sprint.
func (r *UserStoryRepository) AssignUserStoryToSprint(id uint, sprintID uint) error {
	return r.DB.Model(&models.UserStory{}).Where("id = ?", id).Update("sprint_id", sprintID).Error
}

// DeleteUserStory removes a user story from the database by its ID.
func (r *UserStoryRepository) DeleteUserStory(id uint) error {
	return r.DB.Delete(&models.UserStory{}, id).Error
}

// GetUserStoryIDsByProjectID retrieves the IDs of all user stories for a given project.
func (r *UserStoryRepository) GetUserStoryIDsByProjectID(tx *gorm.DB, projectID uint) ([]uint, error) {
	var ids []uint
	err := tx.Model(&models.UserStory{}).Where("project_id = ?", projectID).Pluck("id", &ids).Error
	return ids, err
}

// DeleteUserStoriesByProjectID deletes all user stories associated with a project.
func (r *UserStoryRepository) DeleteUserStoriesByProjectID(tx *gorm.DB, projectID uint) error {
	return tx.Where("project_id = ?", projectID).Delete(&models.UserStory{}).Error
}
