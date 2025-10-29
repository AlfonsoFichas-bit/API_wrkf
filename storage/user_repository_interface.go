package storage

import "github.com/buga/API_wrkf/models"

// UserRepository defines the interface for user data operations.
// This allows for mocking in tests.
type IUserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
    GetNonAdminUsersNotInProject(assignedUserIDs []uint) ([]models.User, error)
}
