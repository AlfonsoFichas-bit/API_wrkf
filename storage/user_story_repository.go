package storage

import (
	"log"

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

// GetUserStoryByID retrieves a single user story by its ID, preloading all related data.
func (r *UserStoryRepository) GetUserStoryByID(id uint) (*models.UserStory, error) {
	var userStory models.UserStory
	err := r.DB.Preload("Project").Preload("CreatedBy").Preload("AssignedTo").First(&userStory, id).Error
	if err == nil {
		if userStory.Points != nil {
			log.Printf("GetUserStoryByID(%d): Fetched user story. Points: %d", id, *userStory.Points)
		} else {
			log.Printf("GetUserStoryByID(%d): Fetched user story. Points is nil.", id)
		}
	}
	return &userStory, err
}

// UpdateUserStory updates an existing user story in the database.
func (r *UserStoryRepository) UpdateUserStory(userStory *models.UserStory) error {
	// Use Select("*") to update all fields, including nil pointers if they were explicitly set to nil.
	// This is a more robust way to handle updates with GORM, especially for nullable fields.
	return r.DB.Model(userStory).Select("*").Updates(userStory).Error
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

// GetUserStoriesBySprintID retrieves all user stories for a given sprint ID.
func (r *UserStoryRepository) GetUserStoriesBySprintID(sprintID uint) ([]models.UserStory, error) {
	var userStories []models.UserStory
	err := r.DB.Where("sprint_id = ?", sprintID).Find(&userStories).Error
	return userStories, err
}
