
package services

import (
	"API_wrkf/models"
	"API_wrkf/storage"
	"errors"
	"time"

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

func (s *UserService) CreateUser(user *models.User) error {
	// Hashear la contraseña antes de guardarla
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Contraseña), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Contraseña = string(hashedPassword)

	return s.Repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

// Login verifies user credentials and returns a JWT token if they are valid.
func (s *UserService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(password))
	if err != nil {
		// If passwords don't match, return the same generic error
		return "", errors.New("invalid credentials")
	}

	// --- Generate JWT Token ---
	claims := jwt.MapClaims{
		"sub": user.ID,
		"nam": user.Nombre,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key from the service struct
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
