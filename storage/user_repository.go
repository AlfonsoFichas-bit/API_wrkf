package storage

import (
	"github.com/buga/API_wrkf/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("correo = ?", email).First(&user).Error
	return &user, err
}

// GetAllUsers retrieves all users from the database.
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

// UpdateUser saves the changes of a user model to the database.
func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

// DeleteUser removes a user from the database by their ID.
func (r *UserRepository) DeleteUser(id uint) error {
	result := r.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Or a custom error indicating user not found
	}
	return nil
}

// GetNonAdminUsersNotInProject retrieves all users who are not admins and are not in the provided list of IDs.
func (r *UserRepository) GetNonAdminUsersNotInProject(assignedUserIDs []uint) ([]models.User, error) {
	var users []models.User
	db := r.DB.Where("role <> ?", models.RoleAdmin)

	if len(assignedUserIDs) > 0 {
		db = db.Where("id NOT IN ?", assignedUserIDs)
	}

	err := db.Find(&users).Error
	return users, err
}
