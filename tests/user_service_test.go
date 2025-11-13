package tests

import (
	"testing"

	"github.com/buga/API_wrkf/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	userService := testApp.UserService

	t.Run("Create and Get User", func(t *testing.T) {
		user := &models.User{
			Nombre:     "John",
			Correo:     "john.doe@test.com",
			Contraseña: "password123",
			Role:       "user",
		}
		err := userService.CreateUser(user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)

		foundUser, err := userService.GetUserByID(user.ID)
		require.NoError(t, err)
		assert.Equal(t, "john.doe@test.com", foundUser.Correo)

		// Check that password is not returned
		assert.Empty(t, foundUser.Contraseña)
	})

	t.Run("Get User By Email", func(t *testing.T) {
		email := "jane.doe@test.com"
		user := &models.User{
			Nombre:     "Jane",
			Correo:     email,
			Contraseña: "password123",
			Role:       "user",
		}
		err := userService.CreateUser(user)
		require.NoError(t, err)

		foundUser, err := userService.GetUserByEmail(email)
		require.NoError(t, err)
		assert.Equal(t, "Jane", foundUser.Nombre)
	})

	t.Run("Password Hashing and Verification via Login", func(t *testing.T) {
		rawPassword := "very-secret-password"
		email := "secret.user@test.com"
		user := &models.User{
			Nombre:     "Secret",
			Correo:     email,
			Contraseña: rawPassword,
			Role:       "user",
		}
		err := userService.CreateUser(user)
		require.NoError(t, err)

		// Verify password by attempting to log in
		_, err = userService.Login(email, rawPassword)
		assert.NoError(t, err, "Login with correct password should succeed")

		_, err = userService.Login(email, "wrong-password")
		assert.Error(t, err, "Login with incorrect password should fail")
	})
}
