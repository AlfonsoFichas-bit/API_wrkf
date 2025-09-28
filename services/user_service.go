package services

import (
	"errors"
	"time"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo      *storage.UserRepository
	jwtSecret []byte
}

func NewUserService(repo *storage.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		Repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

// hashPassword is a helper function to hash passwords.
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CreateAdminUser creates a new user with the admin platform role.
func (s *UserService) CreateAdminUser(user *models.User) error {
	hashedPass, err := hashPassword(user.Contraseña)
	if err != nil {
		return err
	}
	user.Contraseña = hashedPass
	user.Role = string(models.RoleAdmin) // Use the constant for admin
	return s.Repo.CreateUser(user)
}

// CreateUser creates a new user with the default 'user' platform role.
func (s *UserService) CreateUser(user *models.User) error {
	hashedPass, err := hashPassword(user.Contraseña)
	if err != nil {
		return err
	}
	user.Contraseña = hashedPass
	// The role will be set to the default 'user' by the database.
	return s.Repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.Repo.GetUserByEmail(email)
}

// GetAllUsers retrieves all users.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

// DeleteUser deletes a user.
func (s *UserService) DeleteUser(id uint) error {
	return s.Repo.DeleteUser(id)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"nam": user.Nombre,
		"rol": user.Role, // This now correctly reflects the platform role
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
