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
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	user.Contraseña = "" // Clear password for security
	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	user.Contraseña = "" // Clear password for security
	return user, nil
}

// GetAllUsers retrieves all users from the service layer.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Clear passwords before returning
	for i := range users {
		users[i].Contraseña = ""
	}

	return users, nil
}

// UpdateUser handles the logic for updating a user's details.
func (s *UserService) UpdateUser(id uint, updatedData *models.User) (*models.User, error) {
	// Retrieve the existing user
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, err // User not found
	}

	// Update fields if new values are provided
	if updatedData.Nombre != "" {
		user.Nombre = updatedData.Nombre
	}
	if updatedData.ApellidoPaterno != "" {
		user.ApellidoPaterno = updatedData.ApellidoPaterno
	}
	if updatedData.ApellidoMaterno != "" {
		user.ApellidoMaterno = updatedData.ApellidoMaterno
	}
	if updatedData.Correo != "" {
		user.Correo = updatedData.Correo
	}

	// If a new password is provided, hash it and update it
	if updatedData.Contraseña != "" {
		hashedPass, err := hashPassword(updatedData.Contraseña)
		if err != nil {
			return nil, err
		}
		user.Contraseña = hashedPass
	}

	// Save the updated user to the database
	if err := s.Repo.UpdateUser(user); err != nil {
		return nil, err
	}

	// Clear password before returning the user
	user.Contraseña = ""
	return user, nil
}

// DeleteUser handles the logic for deleting a user by their ID.
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

// GenerateJWT generates a JWT token for a given user ID, including their role.
func (s *UserService) GenerateJWT(userID uint) (string, error) {
	// Fetch the user to get their role
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"nam": user.Nombre,
		"rol": user.Role, // Add the user's role to the token claims
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
