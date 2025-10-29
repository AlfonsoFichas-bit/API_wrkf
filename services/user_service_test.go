package services

import (
	"errors"
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/buga/API_wrkf/storage/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Login(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	userService := NewUserService(mockRepo, "test_secret")

	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	mockUser := &models.User{
		ID:          1,
		Nombre:      "Test",
		Correo:      "test@example.com",
		Contraseña:  string(hashedPassword),
		Role:        "user",
	}

	// --- Test Case 1: Successful Login ---
	t.Run("Successful Login", func(t *testing.T) {
		// Setup mock expectation
		mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil).Once()

		// Call the function
		token, err := userService.Login("test@example.com", password)

		// Assertions
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t) // Verify that the mock was called as expected
	})

	// --- Test Case 2: Failed Login - Invalid Credentials ---
	t.Run("Failed Login - Invalid Credentials", func(t *testing.T) {
		// Setup mock expectation
		mockRepo.On("GetUserByEmail", "test@example.com").Return(mockUser, nil).Once()

		// Call the function with the wrong password
		token, err := userService.Login("test@example.com", "wrongpassword")

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

    // --- Test Case 3: Failed Login - User Not Found ---
    t.Run("Failed Login - User Not Found", func(t *testing.T) {
        // Setup mock expectation
        mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, errors.New("record not found")).Once()

        // Call the function
        token, err := userService.Login("notfound@example.com", "password123")

        // Assertions
        assert.Error(t, err)
        assert.Equal(t, "invalid credentials", err.Error())
        assert.Empty(t, token)
        mockRepo.AssertExpectations(t)
    })
}
